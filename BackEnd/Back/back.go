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
	// action types of cache note requestLogin
	Save = 1
	Del  = 2
	Get  = 3
	Edit = 4
	// action types of cache Login requestLogin
	Login  = 1
	signUp = 2
	// responseNote types of responseNote
	successful    = 0
	illegalAccess = 1
	noNote        = 2
	noAccess      = 3
	noUserName    = 4
	wrongPass     = 5
	userNameExist = 6
)

var requestTypeMap = map[string]int{
	http.MethodGet:    Get,
	http.MethodPut:    Edit,
	http.MethodPost:   Save,
	http.MethodDelete: Del,
}

type requestNote struct {
	title string `json:"title"`
	text  string `json:"text"`
}
type requestLogin struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type responseLogin struct {
	jwt       string         `json:"jwt"`
	Name      string         `json:"name"`
	Notes     []responseNote `json:"notes"`
	missCache bool           `json:"misscache"`
}

type responseNote struct {
	note      string `json:"note"`
	title     string `json:"title"`
	noteId    string `json:"noteid"`
	missCache bool   `json:"misscache"`
}

//todo config file
var jwtTries map[string]int = map[string]int{}
var jwtTime = map[string]time.Time{}

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
		if jwtTries[tokenString] <= minuteTryLimit {
			jwtTries[tokenString] += 1
			return claims["authorId"].(string)
		} else {
			// try limit reached
			return "l"
		}
	} else {
		fmt.Println(err)
		return ""
	}
}
func main() {
	preLoad()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//check token
		loginToken := r.Header.Get("jwt")
		jwt := loginToken
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
			} else if authorId == "l" {
				//try limit reached
				w.WriteHeader(http.StatusTooManyRequests)
			}
			//jwt real
		}
		//extract front requestLogin
		noteId, note, noteTitle, done := extractRequest(w, r)
		if done {
			return
		}
		//Get data from cache
		cRes := cache_client.RequestNoteCache(requestTypeMap[r.Method], note, noteTitle, noteId, authorId)
		//handle req and Get res
		res, handleErr := handleNoteRequest(w, r, cRes)
		if handleErr {
			return
		}
		//send responseNote to front
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

func preLoad() {
	ticker := time.NewTicker(time.Minute * time.Duration(20))
	go func(ticker *time.Ticker) {
		for range ticker.C {
			for ts, t := range jwtTime {
				if time.Now().Add(time.Minute*time.Duration(20*2)).Unix() < t.Unix() {
					//jwt is expired
					delete(jwtTime, ts)
					delete(jwtTries, ts)
				}
			}
		}
	}(ticker)
	cache_client.Connect()
	client = cache_client.C
	contextVar = cache_client.Ctx
	//loginRes := requestLoginCache(signUp, "amir123", "Xamir266")
	//fmt.Println(loginRes.Exist)
	//fmt.Println(loginRes.WrongPass)
	//fmt.Println(loginRes.UserId)
}

func handleLoginRequest(w http.ResponseWriter, r *http.Request) {
	loginJson := getRequestBody(r)
	var loginData requestLogin
	err := json.Unmarshal([]byte(loginJson), &loginData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var ActionType int
	if r.Method == http.MethodPost {
		ActionType = 1
	} else if r.Method == http.MethodPut {
		ActionType = 2
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	cRes := cache_client.RequestLoginCache(ActionType, loginData.Name, loginData.Email, loginData.Password)
	var res responseLogin
	if ActionType == Login {
		if cRes.Exist {
			if cRes.WrongPass {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusAccepted)
				//todo config session length
				jwt := createJWT(20, cRes.UserId)
				res.jwt = jwt
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	} else if ActionType == signUp {
		if cRes.Exist {
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			jwt := createJWT(20, cRes.UserId)
			res.jwt = jwt
			w.WriteHeader(http.StatusCreated)
		}
	}
	resJson, _ := json.Marshal(res)
	_, err = w.Write(resJson)
	if err != nil {
		return
	}
	return
}

func extractRequest(w http.ResponseWriter, r *http.Request) (string, string, string, bool) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return "", "", "", true
	}
	var noteId string
	urlList := strings.Split(r.URL.Path, "/")
	if len(urlList) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", "", true
	} else if len(urlList) == 2 {
		noteId = urlList[1]
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", "", true
	}
	var noteJson string
	noteJson = getRequestBody(r)
	var noteObj requestNote
	err := json.Unmarshal([]byte(noteJson), &requestNote{})
	if err != nil {
		return "", "", "", false
	}
	return noteId, noteObj.text, noteObj.title, false
}

func handleNoteRequest(w http.ResponseWriter, r *http.Request, cRes *pb.CacheNoteResponse) (responseNote, bool) {
	var res responseNote
	switch r.Method {
	case http.MethodGet:
		// Get the note.
		res.missCache = cRes.MissCache
		if cRes.Access {
			if cRes.Exist {
				res.note = cRes.Note
				w.WriteHeader(http.StatusFound)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case http.MethodPost:
		// Create a new note.
		w.WriteHeader(http.StatusAccepted)
		res.noteId = cRes.NoteId
	case http.MethodPut:
		// Update an existing note.
		if cRes.Access {
			if cRes.Exist {
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case http.MethodDelete:
		// Remove the note.
		if cRes.Access {
			if cRes.Exist {
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return responseNote{}, true
	}
	return res, false
}

func getRequestBody(r *http.Request) string {
	b, _ := ioutil.ReadAll(r.Body)
	bod := string(b)
	return bod
}
