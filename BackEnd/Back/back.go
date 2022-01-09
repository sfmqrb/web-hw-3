package main

import (
	"encoding/json"
	"fmt"
	"hw3/BackEnd/cache_client"
	pb "hw3/BackEnd/cacheproto"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	// action types of cache text requestLogin
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
	Title string `json:"title"`
	Text  string `json:"text"`
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
	Text      string `json:"text"`
	Title     string `json:"title"`
	NoteId    string `json:"_id"`
	MissCache bool   `json:"misscache"`
}
type Config struct {
	Port           string `json:"port"`
	SessionLimit   int    `json:"sessionLimit"`
	MinuteTryLimit int    `json:"minuteTryLimit"`
}

var config Config
var jwtTries map[string]int = map[string]int{}
var jwtTime = map[string]time.Time{}

var hmacSampleSecret = []byte("my_secret_key")

func toMyNote(notes []*pb.Note) []responseNote {
	var pbNotes []responseNote
	pbNotes = make([]responseNote, len(notes))
	for i := 0; i < len(notes); i++ {
		pbNotes[i] = responseNote{
			Text:      notes[i].Text,
			Title:     notes[i].Title,
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
			fmt.Println("token tick " + ts)
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
	fmt.Println(r.Method)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization, jwt")
	fmt.Println(r)
	fmt.Println(r.Body)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	//fmt.Println(r.Method)
	//check token
	loginToken := r.Header.Get("jwt")
	fmt.Println(loginToken)
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
	//extract front requestLogin
	noteId, note, noteTitle, done := extractRequest(w, r)
	fmt.Println("NoteId:", noteId)
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
	fmt.Println(resJson)
	_, err := w.Write(resJson)
	if err != nil {
		return
	}
}
func main() {
	preLoad()
	//headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	//originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	//methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	//cors := handlers.CORS(
	//	handlers.AllowedHeaders([]string{"content-type"}),
	//	handlers.AllowedOrigins([]string{"*"}),
	//	handlers.AllowCredentials(),
	//)
	//router := mux.NewRouter()
	//router.HandleFunc("/signup", ac.SignUp).Methods("POST")
	//router.HandleFunc("/signin", ac.SignIn).Methods("POST")
	http.HandleFunc("/", HandleRequest)
	//router.Use(cors)
	cache_client.Connect()
	log.Fatal(http.ListenAndServe(":"+config.Port, nil))
	//http.HandleFunc("/", )
	//if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
	//	log.Fatal(err)
	//}
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
	//todo endpoint
	urlList := strings.Split(r.URL.Path, "/")
	var ActionType int
	if urlList[1] == "users" {
		ActionType = 2
	} else if urlList[1] == "auth" {
		ActionType = 1
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cRes := cache_client.RequestLoginCache(ActionType, loginData.UserName, loginData.Name, loginData.Password)
	if cRes == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var res responseLogin
	if ActionType == Login {
		if cRes.Exist {
			if cRes.WrongPass {
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				w.WriteHeader(http.StatusAccepted)
				//todo Config session length
				jwt := createJWT(20, cRes.UserId)
				res.Jwt = jwt
				res.Notes = toMyNote(cRes.Notes)
				res.Name = cRes.Name
			}
		} else {
			w.WriteHeader(http.StatusNotAcceptable)
		}
	} else if ActionType == signUp {
		if cRes.Exist {
			w.WriteHeader(http.StatusNotAcceptable)
		} else {
			jwt := createJWT(20, cRes.UserId)
			res.Jwt = jwt
			res.Name = cRes.UserName
			w.WriteHeader(http.StatusCreated)
		}
	}
	resJson, _ := json.Marshal(res)
	fmt.Println(res)
	b, errw := w.Write(resJson)
	fmt.Println(b)
	if errw != nil {
		return
	}
	return
}

func extractRequest(w http.ResponseWriter, r *http.Request) (string, string, string, bool) {
	//if r.URL.Path != "/" {
	//	http.NotFound(w, r)
	//	return "", "", "", true
	//}
	var noteId string
	urlList := strings.Split(r.URL.Path, "/")
	if len(urlList) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", "", true
	} else if len(urlList) == 3 {
		fmt.Println(urlList)
		fmt.Println(len(urlList))
		noteId = urlList[2]
		fmt.Println(noteId)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		return "", "", "", true
	}
	var noteJson string
	noteJson = getRequestBody(r)

	fmt.Println("notJson "+ noteJson)
	var noteObj requestNote
	err := json.Unmarshal([]byte(noteJson), &noteObj)
	if err != nil {
		return noteId, "", "", false
	}
	return noteId, noteObj.Text, noteObj.Title, false
}

func handleNoteRequest(w http.ResponseWriter, r *http.Request, cRes *pb.CacheNoteResponse) (responseNote, bool) {
	var res responseNote
	fmt.Println(cRes.Access)
	fmt.Println(cRes.Exist)
	switch r.Method {
	case http.MethodGet:
		// Get the text.
		res.MissCache = cRes.MissCache
		if cRes.Access {
			if cRes.Exist {
				res.Text = cRes.Note

				res.NoteId = cRes.NoteId
				w.WriteHeader(http.StatusOK)
			} else {
				w.WriteHeader(http.StatusNoContent)
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	case http.MethodPost:
		// Create a new text.
		w.WriteHeader(http.StatusAccepted)
		res.NoteId = cRes.NoteId
	case http.MethodPut:
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
	case http.MethodDelete:
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
