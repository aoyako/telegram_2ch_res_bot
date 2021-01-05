package dvach

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/telegram"
)

// APIWorkerDvach represents struct to work with external api
type APIWorkerDvach struct {
	cnt        *controller.Controller
	Sender     telegram.Sender
	RequestURL *RequestURL
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
func NewAPIWorkerDvach(cnt *controller.Controller, snd telegram.Sender, rU *RequestURL) *APIWorkerDvach {
	return &APIWorkerDvach{
		cnt:        cnt,
		Sender:     snd,
		RequestURL: rU,
	}
}

// InitiateSending loads data from server and sending it to users
func (dw *APIWorkerDvach) InitiateSending() {
	fmt.Println("started sending")

	boardSubs := make(map[string][]logic.Publication)
	subs := dw.cnt.GetAllSubs()

	for i := range subs {
		boardSubs[subs[i].Board] = append(boardSubs[subs[i].Board], subs[i])
	}

	for key := range boardSubs {
		dw.processBoard(boardSubs[key], key)
	}
}

// Process request from board
func (dw *APIWorkerDvach) processBoard(subs []logic.Publication, board string) {
	resp, err := http.Get(fmt.Sprintf(dw.RequestURL.AllThreadsURL, board))
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

	users := make([][]logic.User, len(subs))
	for subID := range subs {
		userToAppend, _ := dw.cnt.GetUsersByPublication(&subs[subID])
		users[subID] = userToAppend
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
					for userID := range users[subID] {
						usedThreads[threadID] = append(usedThreads[threadID], UserRequest{
							User:    &users[subID][userID],
							Request: subTypes[subID],
						})
					}
					break
				}
			}
		}
	}

	for threadID, subsList := range usedThreads {
		URLThreadID := list.Threads[threadID].ID
		dw.processThread(board, URLThreadID, subsList)
	}
}

// Process requests from thread
func (dw *APIWorkerDvach) processThread(board, URLThreadID string, subsList []UserRequest) {
	resp, err := http.Get(fmt.Sprintf(dw.RequestURL.ThreadURL, board, URLThreadID))
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
				go dw.Sender.Send(fileReceivers, fmt.Sprintf(dw.RequestURL.ResourceURL, file.Path), "")
			}

			if post.Timestamp > currentTimestamp {
				currentTimestamp = post.Timestamp
			}
		}
	}

	dw.cnt.SetLastTimestamp(currentTimestamp)
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
