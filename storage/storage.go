package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// User interface defines methods for User Storage
type User interface {
	Register(user *logic.User) error              // Adds user in databse
	Unregister(user *logic.User) error            // Removes user from database
	GetByChatID(chatID uint) (*logic.User, error) // Returns user by chat id
}

// Subscription interface defines methods for User Storage
type Subscription interface {
	Add(user *logic.User, publication *logic.Publication) error    // Adds new subscription to user with publication
	Remove(user *logic.User, publication *logic.Publication) error // Removes existing sybscription from user
	Update(user *logic.User, publication *logic.Publication) error // Updates selected subscription
	GetByUser(user *logic.User) ([]logic.Publication, error)       // Returns list of user's subscriptions
}

// Info interface definces methods for Info Storage
type Info interface {
	GetLastTimestamp() uint    // Returns time of the latest post
	SetLastTimestamp(tsp uint) // Sets time of the latest post
}

// Storage struct is used to access database
type Storage struct {
	User
	Subscription
	Info
}

// NewStorage constructor of Storage
func NewStorage(db *gorm.DB) *Storage {
	return &Storage{
		User:         NewUserPostgres(db),
		Subscription: NewSubscriptionPostgres(db),
		Info:         NewInfoPostgres(db),
	}
}
