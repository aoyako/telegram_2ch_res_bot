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
	return nil
}

// Unregister removes user from database
func (userStorage *UserPostgres) Unregister(user *logic.User) error {
	return nil
}
