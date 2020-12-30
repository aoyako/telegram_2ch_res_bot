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
func (ucon *UserController) Register(chatID uint64) error {
	user := &logic.User{
		ChatID: chatID,
	}
	return ucon.stg.Register(user)
}

// Unregister performs user deregistration
func (ucon *UserController) Unregister(chatID uint64) error {
	user := &logic.User{
		ChatID: chatID,
	}
	return ucon.stg.Unregister(user)
}

// GetUserByPublication returns owner of publication
func (ucon *UserController) GetUserByPublication(pub *logic.Publication) (*logic.User, error) {
	user, err := ucon.stg.GetUserByPublication(pub)

	return user, err
}
