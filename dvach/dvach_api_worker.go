package dvach

import (
	"log"
	"regexp"
	"strings"

	"github.com/aoyako/telegram_2ch_res_bot/controller"
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/telegram"
)

// APIWorkerDvach represents struct to work with external api
type APIWorkerDvach struct {
	cnt       *controller.Controller
	Sender    telegram.Sender
	Requester Requester
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
func NewAPIWorkerDvach(cnt *controller.Controller, snd telegram.Sender, req Requester) *APIWorkerDvach {
	return &APIWorkerDvach{
		cnt:       cnt,
		Sender:    snd,
		Requester: req,
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

	lastTimestamp := dw.cnt.GetLastTimestamp()
	for key := range boardSubs {
		go dw.processBoard(boardSubs[key], key, lastTimestamp, boardWaiter)
	}

	var lastReceivedTimestamp uint64
	for i := 0; i < len(boardSubs); i++ {
		tmp := <-boardWaiter
		if tmp > lastReceivedTimestamp {
			lastReceivedTimestamp = tmp
		}
	}
	dw.cnt.SetLastTimestamp(lastReceivedTimestamp)
}

// Process request from board
func (dw *APIWorkerDvach) processBoard(subs []logic.Publication, board string, lastTimestamp uint64, waiter chan uint64) {
	list := dw.Requester.GetAllThreads(board)

	users := make([][]logic.User, len(subs))
	for subID := range subs {
		userToAppend, _ := dw.cnt.GetUsersByPublication(&subs[subID])
		users[subID] = userToAppend
	}

	usedThreads := make(map[int]([]UserRequest))

	subValidator := make([]func(string) bool, len(subs))
	subTypes := make([]SourceType, len(subs))
	for i, sub := range subs {
		subValidator[i] = ParseKeywords(sub.Tags)
		subTypes[i] = ParseTypes(sub.Type)
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
		go dw.processThread(board, URLThreadID, subsList, lastTimestamp, threadWaiter)
	}

	var lastReceivedTimestamp uint64
	for i := 0; i < len(usedThreads); i++ {
		tmp := <-threadWaiter
		if tmp > lastReceivedTimestamp {
			lastReceivedTimestamp = tmp
		}
	}

	waiter <- lastReceivedTimestamp
}

// Process requests from thread
func (dw *APIWorkerDvach) processThread(board, URLThreadID string, subsList []UserRequest, lastTimestamp uint64, waiter chan uint64) {
	threadData := dw.Requester.GetThread(board, URLThreadID)
	currentTimestamp := lastTimestamp
	for _, post := range threadData.ThreadPosts[0].Posts {
		if post.Timestamp > lastTimestamp {
			files := post.Files
			for _, file := range files {
				fileReceivers := make([]*logic.User, 0)
				for subID := range subsList {
					if CheckFileExtension(file.Name, subsList[subID].Request) {
						fileReceivers = append(fileReceivers, subsList[subID].User)
					}
				}
				go dw.Sender.Send(fileReceivers, dw.Requester.GetResourceURL(file.Path), URLThreadID)
			}

			if post.Timestamp > currentTimestamp {
				currentTimestamp = post.Timestamp
			}
		}
	}

	waiter <- currentTimestamp
}

// CheckFileExtension returns true if filename is user's selected type
func CheckFileExtension(filename string, req SourceType) bool {
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

// ParseKeywords retruns function to validate keywords
func ParseKeywords(s string) func(string) bool {
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

// ParseTypes returns types from s as [.img.gif.webm]
func ParseTypes(s string) SourceType {
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
