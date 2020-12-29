package controller

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// User interface defines methods for User Controller
type User interface {
	Register(chatID uint64) error   // Performs user registration
	Unregister(chatID uint64) error // Performs user deregistration
}

// Subscription interface defines methods for Publication Controller
type Subscription interface {
	Add(chatID uint64, request string) error                    // Adds new subscription to user with publication
	Remove(chatID uint64, number uint) error                    // Removes existing sybscription from user
	Update(chatID uint64, request string) error                 // Updates selected subscription
	GetSubsByChatID(chatID uint64) ([]logic.Publication, error) // Returns all user's subs
}

// Info interface definces methods for Info Controller
type Info interface {
	GetLastTimestamp() uint64    // Returns time of the latest post
	SetLastTimestamp(tsp uint64) // Sets time of the latest post
}

// Controller struct is used to access database
type Controller struct {
	User
	Subscription
	Info
}

// NewController constructor of Controller
func NewController(stg *storage.Storage) *Controller {
	return &Controller{
		User:         NewUserController(stg),
		Subscription: NewSubscriptionController(stg),
		Info:         NewInfoController(stg),
	}
}
