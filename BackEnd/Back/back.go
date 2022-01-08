package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"hw3/BackEnd/cache_client"
	pb "hw3/BackEnd/cacheproto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// action types of cache note request
	save = 1
	del  = 2
	get  = 3
	// action types of cache Login request
	Login  = 1
	signUp = 2
	// response types of response
	successful    = 0
	illegalAccess = 1
	noNote        = 2
	noAccess      = 3
	noUserName    = 4
	wrongPass     = 5
	userNameExist = 6
)

var requestTypeMap = map[string]int{
	http.MethodGet:    get,
	http.MethodPut:    save,
	http.MethodPost:   save,
	http.MethodDelete: del,
}

var tokenIdMap = map[string]string{
	"123456": "1",
}

type loginRequest struct {
	ActionType int    `json:"type"`
	UserName   string `json:"user"`
	Password   string `json:"pass"`
}
type response struct {
	responseType int
	note         string
	noteId       string
	jwt          string
	missCache    bool
}

//todo config file
var jwtTries map[string]int = map[string]int{}

var minuteTryLimit int = 10
var hmacSampleSecret = []byte("my_secret_key")
var client = cache_client.C
var contextVar = cache_client.Ctx

func createJWT(sessionLength int, authorId string) string {
	now := time.Now()
	until := now.Add(time.Minute * time.Duration(sessionLength))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"nbf":      now.Unix(),
		"exp":      until.Unix(),
		"authorId": authorId,
	})
	tokenString, _ := token.SignedString(hmacSampleSecret)
	jwtTries[tokenString] += 1
	ticker := time.NewTicker(time.Minute)
	go func(ts string, ticker *time.Ticker) {
		for range ticker.C {
			jwtTries[ts] = 0
		}
	}(tokenString, ticker)
	return tokenString
}
func verifyJWT(tokenString string) string {
	token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSampleSecret, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["authorId"].(string)
	} else {
		fmt.Println(err)
		return ""
	}
}
func main() {
	ticker := time.NewTicker(time.Minute * time.Duration(20))
	go func(ticker *time.Ticker) {
		for range ticker.C {
			for ts := range jwtTries {
				jwtTries[ts] = 0
			}
		}
	}(ticker)
	cache_client.Connect()
	client = cache_client.C
	contextVar = cache_client.Ctx
	loginRes := requestLoginCache(signUp, "amir123", "Xamir266")
	fmt.Println(loginRes.Exist)
	fmt.Println(loginRes.WrongPass)
	fmt.Println(loginRes.UserId)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//check token
		loginToken := r.Header.Get("author")
		jwt := tokenIdMap[loginToken]
		var authorId string
		if jwt == "" {
			//Login or create user
			handleLoginRequest(w, r)
			return
		} else {
			//check jwt
			authorId = verifyJWT(jwt)
			if authorId == "" {
				//jwt unreal
				w.WriteHeader(http.StatusNonAuthoritativeInfo)
				return
			}
			//jwt real
		}
		//extract front request
		noteId, note, done := extractRequest(w, r)
		if done {
			return
		}
		//get data from cache
		cRes := requestNoteCache(requestTypeMap[r.Method], note, noteId, authorId)
		//handle req and get res
		res, handleErr := handleNoteRequest(w, r, cRes)
		if handleErr {
			return
		}
		//send response to front
		resJson, _ := json.Marshal(res)
		_, err := w.Write(resJson)
		if err != nil {
			return
		}
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
	cache_client.Connect()
}

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	loginJson := getRequestBody(r)
	var loginData loginRequest
	err := json.Unmarshal([]byte(loginJson), &loginData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cRes := requestLoginCache(loginData.ActionType, loginData.UserName, loginData.Password)
	var res response
	if loginData.ActionType == Login || cRes.Exist {
		if cRes.WrongPass {
			res.responseType = wrongPass
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.WriteHeader(http.StatusAccepted)
			//todo config session length
			jwt := createJWT(20, cRes.UserId)
			res.jwt = jwt
			res.responseType = successful
		}
	} else {
		res.responseType = 4
		w.WriteHeader(http.StatusNotAcceptable)

	}
	if loginData.ActionType == signUp || cRes.Exist {
		res.responseType = userNameExist
		w.WriteHeader(http.StatusNotAcceptable)
	} else {
		res.responseType = successful
		jwt := createJWT(20, cRes.UserId)
		res.jwt = jwt
		w.WriteHeader(http.StatusCreated)
	}
	resJson, _ := json.Marshal(res)
	_, err = w.Write(resJson)
	if err != nil {
		return
	}
	return
}

func extractRequest(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return "", "", true
	}
	var noteId string
	urlList := strings.Split(r.URL.Path, "/")
	if len(urlList) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", true
	} else if len(urlList) == 2 {
		noteId = urlList[1]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", true
	}
	var note string
	note = getRequestBody(r)
	return noteId, note, false
}

func handleNoteRequest(w http.ResponseWriter, r *http.Request, cRes *pb.CacheNoteResponse) (response, bool) {
	var res response
	switch r.Method {
	case http.MethodGet:
		// get the note.
		res.missCache = cRes.MissCache
		if cRes.Access {
			if cRes.Exist {
				res.note = cRes.Note
				res.responseType = successful
				w.WriteHeader(http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNoContent)
				res.responseType = noNote
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			res.responseType = illegalAccess
		}
	case http.MethodPost:
		// Create a new note.
		w.WriteHeader(http.StatusAccepted)
		res.responseType = successful
		res.noteId = cRes.NoteId
	case http.MethodPut:
		// Update an existing note.
		if cRes.Access {
			if cRes.Exist {
				res.responseType = successful
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusNoContent)
				res.responseType = noNote
			}
		} else {
			res.responseType = illegalAccess
			w.WriteHeader(http.StatusUnauthorized)
		}
	case http.MethodDelete:
		// Remove the note.
		if cRes.Access {
			if cRes.Exist {
				res.responseType = successful
				w.WriteHeader(http.StatusAccepted)
			} else {
				res.responseType = noNote
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			res.responseType = illegalAccess
			w.WriteHeader(http.StatusUnauthorized)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return response{}, true
	}
	return res, false
}
func requestNoteCache(requestType int, note string, noteId string, authorId string) *pb.CacheNoteResponse {
	cacheNoteResponse, err := client.CacheNoteRPC(contextVar, &pb.CacheNoteRequest{
		RequestType: int32(requestType),
		NoteId:      noteId,
		AuthorId:    authorId,
		Note:        note,
	})
	if err != nil {
		return nil
	}
	return cacheNoteResponse
}
func requestLoginCache(requestType int, userName string, pass string) *pb.CacheLoginResponse {
	cacheLoginResponse, err := client.CacheLoginRPC(contextVar, &pb.CacheLoginRequest{
		RequestType: int32(requestType),
		User:        userName,
		Pass:        pass,
	})
	if err != nil {
		return nil
	}
	return cacheLoginResponse
}

func getRequestBody(r *http.Request) string {
	b, _ := ioutil.ReadAll(r.Body)
	bod := string(b)
	return bod
}
