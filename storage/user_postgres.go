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
	var tUser logic.User
	exists := userStorage.db.Where("chat_id = ?", user.ChatID).First(&tUser)
	if !(exists.RowsAffected > 0) {
		result := userStorage.db.Create(user)
		return result.Error
	}
	return nil
}

// Unregister removes user from database
func (userStorage *UserPostgres) Unregister(user *logic.User) error {
	result := userStorage.db.Delete(user)
	return result.Error
}

// GetUserByChatID returns user by chat id
func (userStorage *UserPostgres) GetUserByChatID(chatID uint64) (*logic.User, error) {
	var user logic.User
	result := userStorage.db.Where("chat_id = ?", chatID).First(&user)
	return &user, result.Error
}

// Update user
func (userStorage *UserPostgres) Update(user *logic.User) error {
	result := userStorage.db.Save(user)
	return result.Error
}
