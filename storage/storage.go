package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// User interface defines methods for User Storage
type User interface {
	Register(user *logic.User) error   // Adds user in databse
	Unregister(user *logic.User) error // Removes user from database
}

// Subscription interface defines methods for User Storage
type Subscription interface {
	Add(user *logic.User, publication *logic.Publication) error    // Adds new subscription to user with publication
	Remove(user *logic.User, publication *logic.Publication) error // Removes existing sybscription from user
	Update(user *logic.User, publication *logic.Publication) error // Updates selected subscription
}

// Storage struct is used to access database
type Storage struct {
	User
	Subscription
}

// NewStorage constructor of Storage
func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		User:         NewUserPostgres(db),
		Subscription: NewSubscriptionPostgres(db),
	}
}
