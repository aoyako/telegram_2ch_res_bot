package controller

import (
	"errors"
	"fmt"
	"regexp"

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

// Add new subscription to user with publication
func (scon *SubscriptionController) Add(chatID uint64, request string) error {
	user, err := scon.stg.GetUserByChatID(chatID)
	if err != nil {
		return err
	}
	user.SubsCount++
	err = scon.stg.User.Update(user)
	if err != nil {
		return err
	}

	publication, err := parseRequest(request)
	if err != nil {
		return err
	}
	publication.UserID = user.ID

	return scon.stg.Subscription.Add(user, publication)
}

// Remove existing sybscription from user
func (scon *SubscriptionController) Remove(chatID uint64, number uint) error {
	user, err := scon.stg.User.GetUserByChatID(chatID)
	if err != nil {
		return fmt.Errorf("Cannot find user with chat_id=%d", chatID)
	}

	subs, err := scon.stg.Subscription.GetSubsByUser(user)
	if err != nil {
		return fmt.Errorf("Cannot get user's subs: %s", err.Error())
	}

	if len(subs) <= int(number) {
		return fmt.Errorf("Number %d extends amount of subscribtions", number)
	}

	return scon.stg.Subscription.Remove(user, &subs[number])
}

// Update selected subscription
// May be used in future updates
func (scon *SubscriptionController) Update(chatID uint64, request string) error {
	return nil
}

// GetSubsByChatID returns all user's subs
func (scon *SubscriptionController) GetSubsByChatID(chatID uint64) ([]logic.Publication, error) {
	user, err := scon.stg.GetUserByChatID(chatID)
	if err != nil {
		return nil, fmt.Errorf("Cannot find user with chat_id=%d", chatID)
	}

	subs, err := scon.stg.GetSubsByUser(user)
	if err != nil {
		return nil, fmt.Errorf("Cannot get user's subs: %s", err.Error())
	}

	return subs, nil
}

// GetAllSubs Returns all publications
func (scon *SubscriptionController) GetAllSubs() []logic.Publication {
	return scon.stg.GetAllSubs()
}

// Parses request string
// Request string format: "board_name {.img | .webm | .gif} ["keyword1",...]"
func parseRequest(req string) (*logic.Publication, error) {
	separator := regexp.MustCompile(` `)
	args := separator.Split(req, 3)
	if len(args) != 3 {
		return nil, errors.New("Bad request")
	}

	return &logic.Publication{
		Board: args[0],
		Tags:  args[2],
		Type:  args[1],
	}, nil
}
