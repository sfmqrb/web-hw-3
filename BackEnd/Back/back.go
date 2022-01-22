package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/sync/singleflight"
	"hw3/BackEnd/cache_client"
	pb "hw3/BackEnd/cacheproto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// action types of cache text requestLogin
	Save   = 1
	Del    = 2
	Get    = 3
	Edit   = 4
	GetAll = 5
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

type requestNote struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Type  string `json:"type"`
}
type requestLogin struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
type responseLogin struct {
	Jwt       string         `json:"jwt"`
	Name      string         `json:"name"`
	Notes     []responseNote `json:"notes"`
	MissCache bool           `json:"misscache"`
}

type responseNote struct {
	Text      string         `json:"text"`
	Title     string         `json:"title"`
	Type      string         `json:"type"`
	NoteId    string         `json:"_id"`
	Notes     []responseNote `json:"notes"`
	MissCache bool           `json:"misscache"`
}
type Config struct {
	Port           string `json:"port"`
	SessionLimit   int    `json:"sessionLimit"`
	MinuteTryLimit int    `json:"minuteTryLimit"`
}

var config Config
var jwtTries = map[string]int{}
var jwtTime = map[string]time.Time{}
var requestGroup singleflight.Group
var hmacSampleSecret = []byte("toooooooooooo secret")

func probNotesToNotes(notes []*pb.Note) []responseNote {
	var pbNotes []responseNote
	pbNotes = make([]responseNote, len(notes))
	for i := 0; i < len(notes); i++ {
		pbNotes[i] = responseNote{
			Text:      notes[i].Text,
			Title:     notes[i].Title,
			Type:      notes[i].Type,
			NoteId:    notes[i].Id,
			MissCache: false,
		}
	}
	return pbNotes
}

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
		if jwtTries[tokenString] <= config.MinuteTryLimit {
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
func HandleRequest(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("new req")
	//fmt.Println(r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, jwt")
	//fmt.Println(r)
	//fmt.Println(r.Body)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
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
			return
		}
		//jwt real
	}
	requestType, noteId, note, noteTitle, noteType, done := extractRequest(w, r)
	if done {
		return
	}
	//Get data from cache
	cRes, e, shared := requestGroup.Do("RequestNoteCache", func() (interface{}, error) {
		return cache_client.RequestNoteCache(requestType, note, noteTitle, noteType, noteId, authorId)
	})
	if shared {
		fmt.Println("shared note request cache")
	}
	//cRes, e := cache_client.RequestNoteCache(requestType, note, noteTitle, noteType, noteId, authorId)
	if e != nil {
		print(e)
	}
	res, handleErr := handleNoteRequest(requestType, w, r, cRes.(*pb.CacheNoteResponse))
	if handleErr {
		return
	}
	resJson, _ := json.Marshal(res)
	_, err := w.Write(resJson)
	if err != nil {
		return
	}
}
func main() {
	preLoad()
	http.HandleFunc("/", HandleRequest)
	cache_client.Connect()
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
	//http.lis
}

func preLoad() {
	file, _ := ioutil.ReadFile("Back/config.json")
	err := json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
		return
	}
	cache_client.Connect()
	ticker := time.NewTicker(time.Minute * time.Duration(config.SessionLimit))
	go func(ticker *time.Ticker) {
		for range ticker.C {
			for ts, t := range jwtTime {
				if time.Now().Add(time.Minute*time.Duration(config.SessionLimit*2)).Unix() < t.Unix() {
					//jwt is expired
					delete(jwtTime, ts)
					delete(jwtTries, ts)
				}
			}
		}
	}(ticker)
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
	urlList := strings.Split(r.URL.Path, "/")
	var ActionType int
	if urlList[1] == "users" {
		ActionType = signUp
	} else if urlList[1] == "auth" {
		ActionType = Login
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cResRaw, e, shared := requestGroup.Do("RequestNoteCache", func() (interface{}, error) {
		return cache_client.RequestLoginCache(ActionType, loginData.UserName, loginData.Name, loginData.Password)
	})
	if shared {
		fmt.Println("shared login request cache")
	}
	if e != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	cRes := cResRaw.(*pb.CacheLoginResponse)
	//cRes := cache_client.RequestLoginCache(ActionType, loginData.UserName, loginData.Name, loginData.Password)
	var res responseLogin
	res.MissCache = cRes.MissCache
	if ActionType == Login {
		if cRes.Exist {
			if cRes.WrongPass {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusAccepted)
				//todo Config session length
				jwt := createJWT(config.SessionLimit, cRes.UserId)
				res.Jwt = jwt
				res.Notes = probNotesToNotes(cRes.Notes)
				res.Name = cRes.Name
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	} else if ActionType == signUp {
		if cRes.Exist {
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			jwt := createJWT(config.SessionLimit, cRes.UserId)
			res.Jwt = jwt
			res.Name = cRes.Name
			w.WriteHeader(http.StatusCreated)
		}
	}
	resJson, _ := json.Marshal(res)
	//fmt.Println(res)
	_, errw := w.Write(resJson)
	//fmt.Println(b)
	if errw != nil {
		return
	}
	return
}

func extractRequest(w http.ResponseWriter, r *http.Request) (int, string, string, string, string, bool) {
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return "", "", "", true
	//}
	var requestType int
	var noteId = "-1"
	urlList := strings.Split(r.URL.Path, "/")
	if len(urlList) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return requestType, "", "", "", "", true
	} else if len(urlList) == 3 {
		noteId = urlList[2]
	}
	//fmt.Println(urlList)
	//fmt.Println(len(urlList))
	//fmt.Println(noteId)
	//} else if len(urlList) == 2 {
	//
	//} else {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return requestType,"", "", "", true
	//}
	switch r.Method {
	case http.MethodPut:
		requestType = Edit
	case http.MethodGet:
		if urlList[1] == "notes" && urlList[2] == "all" {
			requestType = GetAll
		} else {
			requestType = Get
		}
	case http.MethodPost:
		requestType = Save
	case http.MethodDelete:
		requestType = Del
	}
	var noteJson string
	noteJson = getRequestBody(r)

	fmt.Println("notJson " + noteJson)
	var noteObj requestNote
	err := json.Unmarshal([]byte(noteJson), &noteObj)
	if err != nil {
		return requestType, noteId, "", "", "", false
	}
	return requestType, noteId, noteObj.Text, noteObj.Title, noteObj.Type, false
}

func handleNoteRequest(actionType int, w http.ResponseWriter, r *http.Request, cRes *pb.CacheNoteResponse) (responseNote, bool) {
	var res responseNote
	//fmt.Println(cRes.Access)
	//fmt.Println(cRes.Exist)
	switch actionType {
	case Get:
		// Get the text.
		res.MissCache = cRes.MissCache
		if cRes.Access {
			if cRes.Exist {
				res.Text = cRes.Note
				res.Title = cRes.Title
				res.NoteId = cRes.NoteId
				res.Type = cRes.Type
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case GetAll:
		// Get all texts.
		res.MissCache = cRes.MissCache
		if cRes.Access {
			if cRes.Exist {
				res.Notes = probNotesToNotes(cRes.Notes)
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case Save:
		// Create a new text.
		w.WriteHeader(http.StatusAccepted)
		res.NoteId = cRes.NoteId
	case Edit:
		// Update an existing text.
		if cRes.Access {
			if cRes.Exist {
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case Del:
		// Remove the text.
		if cRes.Exist {
			if cRes.Access {
				w.WriteHeader(http.StatusAccepted)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
			}
		} else {
			w.WriteHeader(http.StatusNoContent)
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
