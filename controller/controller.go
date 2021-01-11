package controller

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// User interface defines methods for User Controller
type User interface {
	Register(chatID uint64) error                                       // Performs user registration
	Unregister(chatID uint64) error                                     // Performs user deregistration
	GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) // Returns owner of publication
}

// Subscription interface defines methods for Publication Controller
type Subscription interface {
	AddNew(chatID uint64, request string) error // Adds new subscription to user with publication
	Create(chatID uint64, request string) error
	Remove(chatID uint64, request string) error                 // Removes existing sybscription from user
	Update(chatID uint64, request string) error                 // Updates selected subscription
	GetSubsByChatID(chatID uint64) ([]logic.Publication, error) // Returns all user's subs
	GetAllSubs() []logic.Publication                            // Returns all publications
	GetAllDefaultSubs() []logic.Publication
	RemoveDefault(chatID uint64, request string) error
	Subscribe(chatID uint64, request string) error
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
