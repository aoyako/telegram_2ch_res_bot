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
func (ucon *UserController) Register(user *logic.User) error {
	return nil
}

// Unregister performs user deregistration
func (ucon *UserController) Unregister(user *logic.User) error {
	return nil
}

// Update user
func (ucon *UserController) Update(user *logic.User) error {
	return nil
}
