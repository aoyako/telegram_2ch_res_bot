package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// UserPostgres is an implementation of storage.User
type UserPostgres struct {
	db *gorm.DB
}

// NewUserPostgres constructor of UserPostgres struct
func NewUserPostgres(db *gorm.DB) *UserPostgres {
	return &UserPostgres{
		db: db,
	}
}

// Register adds user in databse
func (userStorage *UserPostgres) Register(user *logic.User) error {
	result := userStorage.db.Create(user)
	return result.Error
}

// Unregister removes user from database
func (userStorage *UserPostgres) Unregister(user *logic.User) error {
	result := userStorage.db.Delete(user)
	return result.Error
}

// GetByChatID returns user by chat id
func (userStorage *UserPostgres) GetByChatID(chatID uint) (*logic.User, error) {
	var user logic.User
	result := userStorage.db.Where("ChatID = ?", chatID).First(&user)
	return &user, result.Error
}
