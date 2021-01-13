package dvach

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// RequestURL describes endpoints of external api
type RequestURL struct {
	AllThreadsURL string
	ThreadURL     string
	ResourceURL   string
}

// Requester gets data from external sources
type Requester interface {
	GetAllThreads(board string) ListResponse
	GetThread(board, threadID string) ThreadData
	GetResourceURL(path string) string
}

// APIRequester gets data from 2ch
type APIRequester struct {
	Requests *RequestURL
}

// NewRequester constructor for APIRequester
func NewRequester(u *RequestURL) *APIRequester {
	return &APIRequester{
		Requests: u,
	}
}

// GetAllThreads returns list of all threads on board
func (r *APIRequester) GetAllThreads(board string) ListResponse {
	resp, err := http.Get(fmt.Sprintf(r.Requests.AllThreadsURL, board))
	if err != nil {
		log.Printf("Error creating request to 2ch: %s", err.Error())
	}

	var list ListResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading request body")
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Printf("Error unmarshalling board request body: %s", err.Error())
	}
	return list
}

// GetThread returns list of posts in the thread with id = threadID
func (r *APIRequester) GetThread(board, threadID string) ThreadData {
	resp, err := http.Get(fmt.Sprintf(r.Requests.ThreadURL, board, threadID))
	if err != nil {
		log.Printf("Error creating request to 2ch.hk: %s", err.Error())
	}

	var threadData ThreadData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading request body")
	}
	err = json.Unmarshal(body, &threadData)
	if err != nil {
		log.Printf("Error unmarshalling thread request body: %s", err.Error())
	}

	return threadData
}

// GetResourceURL converts relative resource path to absolute
func (r *APIRequester) GetResourceURL(path string) string {
	return fmt.Sprintf(r.Requests.ResourceURL, path)
}
