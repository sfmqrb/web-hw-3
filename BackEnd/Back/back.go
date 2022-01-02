package back

import (
	"encoding/json"
	"hw3/BackEnd/cache_client"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	save = 1
	del  = 2
	get  = 3
	//
	successful    = 0
	illegalAccess = 1
	noNote        = 2
	noAccess      = 3
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

type response struct {
	responseType int
	note         string
	noteId       string
	missCache    bool
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//check token
		loginToken := r.Header.Get("author")
		authorId := tokenIdMap[loginToken]
		if authorId == "" {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			return
		}
		//extract front request
		noteId, note, done := extractRequest(w, r)
		if done {
			return
		}
		//get data from cache
		cRes := requestCache(requestTypeMap[r.Method], note, noteId, authorId)
		//handle req and get res
		res, handleErr := handleRequest(w, r, cRes)
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
	note = getRequestNote(r)
	return noteId, note, false
}

func handleRequest(w http.ResponseWriter, r *http.Request, cRes cache_client.CacheResponse) (response, bool) {
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
func requestCache(requestType int, note string, noteId string, authorId string) cache_client.CacheResponse {
	cReq := cache_client.CacheRequest{RequestType: requestType, NoteId: noteId, AuthorId: authorId, Note: note}
	cRes := cache_client.SendRequestCache(cReq)
	return cRes // type cacheResponse
}

func getRequestNote(r *http.Request) string {
	b, _ := ioutil.ReadAll(r.Body)
	bod := string(b)
	return bod
}
