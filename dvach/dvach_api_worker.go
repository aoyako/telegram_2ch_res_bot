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
	log.Println("started sending")

	boardSubs := make(map[string][]logic.Publication)
	subs := dw.cnt.GetAllSubs()

	for i := range subs {
		boardSubs[subs[i].Board] = append(boardSubs[subs[i].Board], subs[i])
	}

	boardWaiter := make(chan uint64, len(boardSubs))

	for key := range boardSubs {
		go dw.processBoard(boardSubs[key], key, boardWaiter)
	}

	var lastTimestamp uint64
	for i := 0; i < len(boardSubs); i++ {
		tmp := <-boardWaiter
		if tmp > lastTimestamp {
			lastTimestamp = tmp
		}
	}
	dw.cnt.SetLastTimestamp(lastTimestamp)
}

// Process request from board
func (dw *APIWorkerDvach) processBoard(subs []logic.Publication, board string, waiter chan uint64) {
	resp, err := http.Get(fmt.Sprintf(dw.RequestURL.AllThreadsURL, board))
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

	users := make([][]logic.User, len(subs))
	for subID := range subs {
		userToAppend, _ := dw.cnt.GetUsersByPublication(&subs[subID])
		users[subID] = userToAppend
	}

	usedThreads := make(map[int]([]UserRequest))

	subValidator := make([]func(string) bool, len(subs))
	subTypes := make([]SourceType, len(subs))
	for i, sub := range subs {
		subValidator[i] = parseKeywords(sub.Tags)
		subTypes[i] = parseTypes(sub.Type)
	}

	for threadID, thread := range list.Threads {
		for subID := range subs {
			if subValidator[subID](thread.Comment) {
				for userID := range users[subID] {
					usedThreads[threadID] = append(usedThreads[threadID], UserRequest{
						User:    &users[subID][userID],
						Request: subTypes[subID],
					})
				}
			}
		}
	}

	threadWaiter := make(chan uint64, len(usedThreads))
	for threadID, subsList := range usedThreads {
		URLThreadID := list.Threads[threadID].ID
		go dw.processThread(board, URLThreadID, subsList, threadWaiter)
	}

	var lastTimestamp uint64
	for i := 0; i < len(usedThreads); i++ {
		tmp := <-threadWaiter
		if tmp > lastTimestamp {
			lastTimestamp = tmp
		}
	}

	waiter <- lastTimestamp
}

// Process requests from thread
func (dw *APIWorkerDvach) processThread(board, URLThreadID string, subsList []UserRequest, waiter chan uint64) {
	resp, err := http.Get(fmt.Sprintf(dw.RequestURL.ThreadURL, board, URLThreadID))
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

	lastTimestamp := dw.cnt.GetLastTimestamp()
	currentTimestamp := lastTimestamp

	for _, post := range threadData.ThreadPosts[0].Posts {
		if post.Timestamp > lastTimestamp {
			files := post.Files
			for _, file := range files {
				fileReceivers := make([]*logic.User, 0)
				for subID := range subsList {
					if checkFileExtension(file.Name, subsList[subID].Request) {
						fileReceivers = append(fileReceivers, subsList[subID].User)
					}
				}
				go dw.Sender.Send(fileReceivers, fmt.Sprintf(dw.RequestURL.ResourceURL, file.Path), URLThreadID)
			}

			if post.Timestamp > currentTimestamp {
				currentTimestamp = post.Timestamp
			}
		}
	}

	waiter <- currentTimestamp
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

// Retruns function to validate keywords
func parseKeywords(s string) func(string) bool {
	d := regexp.MustCompile("\"\\|")
	c := regexp.MustCompile("\"&")
	disjunction := d.Split(s, -1)
	conjunction := make([][]string, len(disjunction))
	negation := make([][]bool, len(disjunction))

	for key := range disjunction {
		conjunction[key] = c.Split(disjunction[key], -1)
		for ckey := range conjunction[key] {
			if strings.HasPrefix(conjunction[key][ckey], "!") {
				negation[key] = append(negation[key], true)
				conjunction[key][ckey] = strings.TrimPrefix(conjunction[key][ckey], "!\"")
			} else {
				negation[key] = append(negation[key], false)
				conjunction[key][ckey] = strings.TrimPrefix(conjunction[key][ckey], "\"")
			}
		}
	}

	conjunction[len(conjunction)-1][len(conjunction[len(conjunction)-1])-1] =
		strings.TrimSuffix(conjunction[len(conjunction)-1][len(conjunction[len(conjunction)-1])-1], "\"")

	return func(input string) bool {
		for dis := range conjunction {
			success := true
			for con := range conjunction[dis] {
				if !(strings.Contains(strings.ToLower(input),
					strings.ToLower(conjunction[dis][con])) != negation[dis][con]) {
					success = false
					break
				}
			}
			if success {
				return true
			}
		}
		return false
	}
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

// Returns if "val" is in sice and it's position
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
