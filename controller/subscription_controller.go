package controller

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// SubscriptionController is an implementation of controller.Subscription
type SubscriptionController struct {
	stg *storage.Storage
}

// NewSubscriptionController constructor of SubscriptionController struct
func NewSubscriptionController(stg *storage.Storage) *SubscriptionController {
	return &SubscriptionController{stg: stg}
}

// AddNew creates a subscription to user with publication
func (scon *SubscriptionController) AddNew(chatID int64, request string) error {
	user, err := scon.stg.GetUserByChatID(chatID)
	if err != nil {
		log.Println("SubscriptionController.AddNew-GetUserByChatID", err)
		return err
	}

	publication, err := parseRequest(request)
	if err != nil {
		log.Println("SubscriptionController.AddNew-parseRequest", err)
		return err
	}

	err = scon.stg.Subscription.Add(user, publication)
	if err != nil {
		log.Println("SubscriptionController.AddNew-Add", err)
		return err
	}

	user.SubsCount++
	err = scon.stg.User.Update(user)
	return err
}

// Create default subscribtion
func (scon *SubscriptionController) Create(chatID int64, request string) error {
	if !scon.stg.IsChatAdmin(chatID) {
		return errors.New("access denied")
	}
	publication, err := parseRequestAlias(request)
	if err != nil {
		log.Println("SubscriptionController.Create-parseRequestAlias", err)
		return err
	}
	publication.IsDefault = true

	return scon.stg.Subscription.AddDefault(publication)
}

// Subscribe user to publication
func (scon *SubscriptionController) Subscribe(chatID int64, request string) error {
	user, err := scon.stg.GetUserByChatID(chatID)
	if err != nil {
		log.Println("SubscriptionController.Subscribe-GetUserByChatID", err)
		return err
	}

	pubID, err := strconv.Atoi(request)
	if err != nil {
		log.Println("SubscriptionController.Subscribe-Atoi", err)
		return errors.New("bad index")
	}
	pubID--

	pubs := scon.stg.GetAllDefaultSubs()
	if pubID >= len(pubs) || pubID < 0 {
		log.Println("SubscriptionController.Subscribe-GetAllDefaultSubs", err)
		return errors.New("bad index")
	}

	err = scon.stg.Subscription.Connect(user, &pubs[pubID])
	if err != nil {
		log.Println("SubscriptionController.Subscribe-Connect", err)
		return err
	}

	user.SubsCount++
	err = scon.stg.User.Update(user)
	return err
}

// Remove existing sybscription from user
func (scon *SubscriptionController) Remove(chatID int64, request string) error {
	user, err := scon.stg.User.GetUserByChatID(chatID)
	if err != nil {
		log.Println("SubscriptionController.Remove-GetUserByChatID", err)
		return fmt.Errorf("cannot find user with chat_id=%d", chatID)
	}

	subs, err := scon.stg.Subscription.GetSubsByUser(user)
	if err != nil {
		log.Println("SubscriptionController.Remove-GetSubsByUser", err)
		return fmt.Errorf("cannot get user's subs: %s", err.Error())
	}

	subID, err := strconv.Atoi(request)
	if err != nil {
		log.Println("SubscriptionController.Remove-Atoi", err)
		return errors.New("bad index")
	}
	subID--

	if subID >= len(subs) || subID < 0 {
		return errors.New("bad index")
	}

	err = scon.stg.Subscription.Disonnect(user, &subs[subID])
	if err != nil {
		log.Println("SubscriptionController.Remove-Disonnect", err)
		return err
	}
	if !subs[subID].IsDefault {
		err = scon.stg.Subscription.Remove(&subs[subID])
		if err != nil {
			log.Println("SubscriptionController.Remove-Remove", err)
			return err
		}
	}

	user.SubsCount--
	err = scon.stg.User.Update(user)
	return err
}

// RemoveDefault deletes default publication
func (scon *SubscriptionController) RemoveDefault(chatID int64, request string) error {
	if !scon.stg.IsChatAdmin(chatID) {
		return errors.New("access denied")
	}

	pubs := scon.stg.Subscription.GetAllDefaultSubs()

	pubID, err := strconv.Atoi(request)
	if err != nil {
		return errors.New("bad index")
	}
	pubID--

	if pubID >= len(pubs) || pubID < 0 {
		return errors.New("bad index")
	}

	users, err := scon.stg.GetUsersByPublication(&pubs[pubID])
	if err != nil {
		return errors.New("bad request")
	}

	for i := range users {
		users[i].SubsCount--
		err := scon.stg.User.Update(&users[i])
		if err != nil {
			log.Println("RemoveDefault.Update User", err)
			return errors.New("bad request")
		}
	}

	return scon.stg.Subscription.Remove(&pubs[pubID])
}

// Update selected subscription
// May be used in future updates
func (scon *SubscriptionController) Update(chatID int64, request string) error {
	return nil
}

// GetSubsByChatID returns all user's subs
func (scon *SubscriptionController) GetSubsByChatID(chatID int64) ([]logic.Publication, error) {
	user, err := scon.stg.GetUserByChatID(chatID)
	if err != nil {
		log.Println("SubscriptionController.GetSubsByChatID-GetUserByChatID", err)
		return nil, fmt.Errorf("cannot find user with chat_id=%d", chatID)
	}

	subs, err := scon.stg.GetSubsByUser(user)
	if err != nil {
		log.Println("SubscriptionController.GetSubsByChatID-GetSubsByUser", err)
		return nil, fmt.Errorf("cannot get user's subs: %s", err.Error())
	}

	return subs, nil
}

// GetAllSubs returns all publications
func (scon *SubscriptionController) GetAllSubs() []logic.Publication {
	return scon.stg.GetAllSubs()
}

// GetAllDefaultSubs returns all default publications
func (scon *SubscriptionController) GetAllDefaultSubs() []logic.Publication {
	return scon.stg.GetAllDefaultSubs()
}

// Parses request string
// Request string format: "board_name {.img | .webm | .gif} "keyword1"[|,&]..."
func parseRequest(req string) (*logic.Publication, error) {
	separator := regexp.MustCompile(` `)
	args := separator.Split(req, 3)
	if len(args) != 3 {
		log.Println("parseRequest - error", args)
		return nil, errors.New("bad request")
	}

	tags := args[2]
	res, err := regexp.MatchString(`^(!?".+"[|&])*!?"[^&|]+"$`, tags)

	if err != nil || !res {
		log.Println("parseRequest - error", args)
		return nil, errors.New("bad request")
	}

	types := args[1]
	res, err = regexp.MatchString(`^(\.[A-Za-z0-9]+)+$`, types)
	if err != nil || !res {
		log.Println("parseRequest - error", args)
		return nil, errors.New("bad request")
	}

	return &logic.Publication{
		Board: args[0],
		Tags:  args[2],
		Type:  args[1],
	}, nil
}

// Parses request string with alias
func parseRequestAlias(req string) (*logic.Publication, error) {
	separator := regexp.MustCompile(` `)
	args := separator.Split(req, 3)

	parser := regexp.MustCompile(`" `)
	words := parser.Split(args[len(args)-1], -1)

	if len(args) != 3 {
		log.Println("parseRequestAlias - error", args)
		return nil, errors.New("bad request")
	}

	tags := strings.TrimSuffix(args[2], " "+words[len(words)-1])
	res, err := regexp.MatchString(`^(!?".+"[|&])*!?"[^&|]+"$`, tags)
	if err != nil || !res {
		log.Println("parseRequestAlias - error", args)
		return nil, errors.New("bad request")
	}

	types := args[1]
	res, err = regexp.MatchString(`^(\.[A-Za-z0-9]+)+$`, types)
	if err != nil || !res {
		log.Println("parseRequestAlias - error", args)
		return nil, errors.New("bad request")
	}

	return &logic.Publication{
		Board: args[0],
		Tags:  tags,
		Type:  types,
		Alias: words[len(words)-1],
	}, nil
}
