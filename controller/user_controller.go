package controller

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// UserController is an implementation of controller.User
type UserController struct {
	stg *storage.Storage
}

// NewUserController constructor of UserController struct
func NewUserController(stg *storage.Storage) *UserController {
	return &UserController{stg: stg}
}

// Register performs user registration
func (ucon *UserController) Register(chatID int64) error {
	user := &logic.User{
		ChatID: chatID,
	}
	return ucon.stg.Register(user)
}

// Unregister performs user deregistration
func (ucon *UserController) Unregister(chatID int64) error {
	user := &logic.User{
		ChatID: chatID,
	}
	return ucon.stg.Unregister(user)
}

// GetUsersByPublication returns subscribers of publication
func (ucon *UserController) GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) {
	users, err := ucon.stg.GetUsersByPublication(pub)

	return users, err
}
