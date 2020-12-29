package controller

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// User interface defines methods for User Controller
type User interface {
	Register(user *logic.User) error   // Performs user registration
	Unregister(user *logic.User) error // Performs user deregistration
	Update(user *logic.User) error     // Updates user
}

// Subscription interface defines methods for Publication Controller
type Subscription interface {
	Add(user *logic.User, publication *logic.Publication) error    // Adds new subscription to user with publication
	Remove(user *logic.User, publication *logic.Publication) error // Removes existing sybscription from user
	Update(user *logic.User, publication *logic.Publication) error // Updates selected subscription
}

// Info interface definces methods for Info Controller
type Info interface {
	GetLastTimestamp() uint    // Returns time of the latest post
	SetLastTimestamp(tsp uint) // Sets time of the latest post
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
