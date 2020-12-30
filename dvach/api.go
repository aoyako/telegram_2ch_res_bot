package dvach

import (
	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
)

// APIWorker for working with external api
type APIWorker interface {
	SetSenderFunc(func(user *logic.User, data interface{}))
	InitiateSending()
}

// APIController for accessing external api
type APIController struct {
	APIWorker
}

// Thread contains thread header
type Thread struct {
	Comment   string  `json:"comment"`
	Lasthit   int64   `json:"lasthit"`
	ID        string  `json:"num"`
	PostCount int     `json:"posts_count"`
	Score     float64 `json:"score"`
	Subject   string  `json:"subject"`
	Timestamp uint64  `json:"timestamp"`
	Views     int     `json:"views"`
}

// ListResponse contains struct to be returned when reading all threads
type ListResponse struct {
	Board   string   `json:"board"`
	Threads []Thread `json:"threads"`
}

// Post contains post data
type Post struct {
	Comment   string `json:"comment"`
	Date      string `json:"date"`
	Timestamp uint64 `json:"timestamp"`
	Files     []File `json:"files"`
}

// File constains file data
type File struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Size int    `json:"size"`
}

// ThreadPost stores info about thread's posts
type ThreadPost struct {
	Posts []Post `json:"posts"`
}

// ThreadData contains every thread data
type ThreadData struct {
	ThreadPosts []ThreadPost `json:"threads"`
}

// NewAPIController constructor of APIController
func NewAPIController(cnt *controller.Controller) *APIController {
	return &APIController{
		APIWorker: NewAPIWorkerDvach(cnt),
	}
}
