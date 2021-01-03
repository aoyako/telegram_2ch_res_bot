package dvach

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/telegram"
)

var DatabaseLock sync.Mutex

// APIWorkerDvach represents struct to work with external api
type APIWorkerDvach struct {
	cnt    *controller.Controller
	Sender telegram.Sender
}

// SourceType specify user's file extensions choice
type SourceType struct {
	Image bool
	Gif   bool
	Webm  bool
}

// UserRequest stores information about user and it's requested file types
// in thread
type UserRequest struct {
	User    *logic.User
	Request SourceType
}

// NewAPIWorkerDvach constructor for APIWorkerDvach
func NewAPIWorkerDvach(cnt *controller.Controller, snd telegram.Sender) *APIWorkerDvach {
	return &APIWorkerDvach{
		cnt:    cnt,
		Sender: snd,
	}
}

// InitiateSending loads data from server and sending it to users
func (dw *APIWorkerDvach) InitiateSending() {
	fmt.Println("started sending")

	boards := make(map[string]bool)
	subs := dw.cnt.GetAllSubs()
	users := make([]*logic.User, len(subs))

	for i := range subs {
		boards[subs[i].Board] = true
		users[i], _ = dw.cnt.GetUserByPublication(&subs[i])
	}

	for board := range boards {
		processBoard(dw, subs, users, board)
	}
}

// Process request from board
func processBoard(dw *APIWorkerDvach, subs []logic.Publication, users []*logic.User, board string) {
	resp, err := http.Get(fmt.Sprintf("https://2ch.hk/%s/threads.json", board))
	if err != nil {
		log.Fatalf("Error creating request to 2ch.hk: %s", err.Error())
	}

	var list ListResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading request body")
	}
	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Fatalf("Error unmarshalling request body: %s", err.Error())
	}

	usedThreads := make(map[int]([]UserRequest))

	subKeywords := make([][]string, len(subs))
	subTypes := make([]SourceType, len(subs))
	for i, sub := range subs {
		subKeywords[i] = parseKeywords(sub.Tags)
		subTypes[i] = parseTypes(sub.Type)
	}

	for threadID, thread := range list.Threads {
		for subID := range subs {
			for _, keyword := range subKeywords[subID] {
				if strings.Contains(thread.Comment, keyword) {
					usedThreads[threadID] = append(usedThreads[threadID], UserRequest{
						User:    users[subID],
						Request: subTypes[subID],
					})
					break
				}
			}
		}
	}

	for threadID, subsList := range usedThreads {
		URLThreadID := list.Threads[threadID].ID
		processThread(dw, board, URLThreadID, subsList)
	}
}

// Process requests from thread
func processThread(dw *APIWorkerDvach, board, URLThreadID string, subsList []UserRequest) {
	resp, err := http.Get(fmt.Sprintf("https://2ch.hk/%s/res/%s.json", board, URLThreadID))
	if err != nil {
		log.Fatalf("Error creating request to 2ch.hk: %s", err.Error())
	}

	var threadData ThreadData
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading request body")
	}
	err = json.Unmarshal(body, &threadData)
	if err != nil {
		log.Fatalf("Error unmarshalling request body: %s", err.Error())
	}

	lastTimestamp := dw.cnt.GetLastTimestamp()
	currentTimestamp := lastTimestamp

	for _, post := range threadData.ThreadPosts[0].Posts {
		if post.Timestamp > lastTimestamp {
			files := post.Files
			for _, file := range files {
				fileReceivers := make([]*logic.User, 0)
				for _, sub := range subsList {
					if checkFileExtension(file.Name, sub.Request) {
						fileReceivers = append(fileReceivers, sub.User)
					}
				}
				go dw.Sender.Send(fileReceivers, fmt.Sprintf("https://2ch.hk%s", file.Path), "")
			}

			if post.Timestamp > currentTimestamp {
				currentTimestamp = post.Timestamp
			}
		}
	}

	DatabaseLock.Lock()
	lastTimestamp = dw.cnt.GetLastTimestamp()
	if lastTimestamp < currentTimestamp {
		dw.cnt.SetLastTimestamp(currentTimestamp)
	}
	DatabaseLock.Unlock()
	fmt.Println("Finished sending")
}

// Returns true if filename is user's selected type
func checkFileExtension(filename string, req SourceType) bool {
	var result bool
	if req.Image {
		result = result || strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".jpg") ||
			strings.HasSuffix(filename, ".jpeg")
	}

	if req.Gif {
		result = result || strings.HasSuffix(filename, ".gif")
	}

	if req.Webm {
		result = result || strings.HasSuffix(filename, ".webm")
	}

	return result
}

// Returns slice of keywords from s as ""keyword1","keyword2",.."
func parseKeywords(s string) []string {
	re := regexp.MustCompile("\",\"")
	res := re.Split(s, -1)
	res[0] = strings.TrimPrefix(res[0], "\"")
	res[len(res)-1] = strings.TrimSuffix(res[len(res)-1], "\"")
	return res
}

// Returns types from s as [.img.gif.webm]
func parseTypes(s string) SourceType {
	var result SourceType

	re := regexp.MustCompile("\\.")
	res := re.Split(s, -1)
	if _, ok := find(res, "img"); ok {
		result.Image = true
	}

	if _, ok := find(res, "gif"); ok {
		result.Gif = true
	}

	if _, ok := find(res, "webm"); ok {
		result.Webm = true
	}

	return result
}

// Returns true and postion of string val in slice "slcie"
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
