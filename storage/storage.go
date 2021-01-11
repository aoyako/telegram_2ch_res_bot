package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// User interface defines methods for User Storage
type User interface {
	Register(user *logic.User) error                                    // Adds user in databse
	Unregister(user *logic.User) error                                  // Removes user from database
	GetUserByChatID(chatID uint64) (*logic.User, error)                 // Returns user by chat id
	Update(user *logic.User) error                                      // Updates user
	GetUserByID(userID uint) (*logic.User, error)                       // Returns user by it's id
	GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) // Returns owner of publication
	IsUserAdmin(user *logic.User) bool
	IsChatAdmin(userID uint64) bool
}

// Subscription interface defines methods for Publicaiton Storage
type Subscription interface {
	Add(user *logic.User, publication *logic.Publication) error       // Adds new subscription to user with publication
	AddDefault(publication *logic.Publication) error                  // Adds new subscription to user with publication
	Remove(publication *logic.Publication) error                      // Removes existing sybscription
	Disonnect(user *logic.User, publication *logic.Publication) error // Disonnect user from publication
	Update(user *logic.User, publication *logic.Publication) error    // Updates selected subscription
	GetSubsByUser(user *logic.User) ([]logic.Publication, error)      // Returns list of user's subscriptions
	GetAllSubs() []logic.Publication                                  // Returns all publications
	GetAllDefaultSubs() []logic.Publication
	Connect(user *logic.User, publication *logic.Publication) error
}

// Info interface definces methods for Info Storage
type Info interface {
	GetLastTimestamp() uint64    // Returns time of the latest post
	SetLastTimestamp(tsp uint64) // Sets time of the latest post
}

// Storage struct is used to access database
type Storage struct {
	User
	Subscription
	Info
}

// NewStorage constructor of Storage
func NewStorage(db *gorm.DB, cfg *InitDatabase) *Storage {
	return &Storage{
		User:         NewUserPostgres(db, cfg),
		Subscription: NewSubscriptionPostgres(db),
		Info:         NewInfoPostgres(db),
	}
}
