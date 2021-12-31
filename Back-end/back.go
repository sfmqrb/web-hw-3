package main

import (
	"encoding/json"
	"fmt"
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

var tokenIdMap = map[string]string{}

type cacheResponse struct {
	note      string
	noteId    string
	exist     bool
	access    bool
	missCache bool
}
type cacheRequest struct {
	requestType int
	noteId      string
	authorId    string
	note        string
}
type response struct {
	responseType int
	note         string
	noteId       string
	missCache    bool
}

type note struct {
	note     string
	authorId string
}

var testMap = map[string]string{

	"a": "b",
}

func main() {
	fmt.Println(testMap["a"])
	x := testMap["b"]
	if x == "" {

		fmt.Println("ok")
	}
	s := strings.Split("adasd,bwe,xcc", ",")
	fmt.Println(len(s))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Printf("Req: %s %s\n", r.Host, r.URL.Path)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			http.NotFound(w, r)
			return
		}
		var noteId string
		urlList := strings.Split(r.URL.Path, "/")
		if len(urlList) == 0 {
			//todo send front
			return
		} else if len(urlList) == 2 {
			noteId = urlList[1]
			fmt.Println(noteId)
		} else {
			//todo
			fmt.Println("bad request")
			return
		}
		loginToken := r.Header.Get("author")
		fmt.Println("Hello World!")
		fmt.Println(getRequestNote(r))
		var note string
		note = getRequestNote(r)
		authorId := tokenIdMap[loginToken]
		cRes := requestCache(requestTypeMap[r.Method], note, noteId, authorId)
		var res response
		switch r.Method {
		case http.MethodGet:
			// get the note.
			res.missCache = cRes.missCache
			if cRes.access {
				if cRes.exist {
					res.note = cRes.note
					res.responseType = successful
				} else {
					res.responseType = noNote
				}
			} else {
				res.responseType = illegalAccess
			}
		case http.MethodPost:
			// Create a new note.
			res.responseType = successful
			res.noteId = cRes.noteId
		case http.MethodPut:
			// Update an existing note.
			if cRes.access {
				if cRes.exist {
					res.responseType = successful
				} else {
					res.responseType = noNote
				}
			} else {
				res.responseType = illegalAccess
			}
		case http.MethodDelete:
			// Remove the note.
			if cRes.access {
				if cRes.exist {
					res.responseType = successful
				} else {
					res.responseType = noNote
				}
			} else {
				res.responseType = illegalAccess
			}
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		resJson, _ := json.Marshal(res)
		_, err := w.Write(resJson)
		if err != nil {
			return
		}
		//handle cache response

	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
func requestCache(requestType int, note string, noteId string, authorId string) cacheResponse {
	cReq := cacheRequest{requestType: requestType, noteId: noteId, authorId: authorId, note: note}
	var cReqJson, err = json.Marshal(cReq)
	if err != nil {
		//todo
		return cacheResponse{}
	}
	fmt.Println(cReqJson)
	//todo cache gRPC
	switch requestType {
	case save:

	case get:

	case del:
	}
	return cacheResponse{} // type cacheResponse
}

func getRequestNote(r *http.Request) string {
	b, _ := ioutil.ReadAll(r.Body)
	bod := string(b)
	return bod
}
