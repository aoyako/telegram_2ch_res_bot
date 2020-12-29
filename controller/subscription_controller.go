package controller

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// SubscriptionController is an implementation of controller.Subscription
type SubscriptionController struct {
	stg *storage.Storage
}

// NewSubscriptionController constructor of SubscriptionController struct
func NewSubscriptionController(stg *storage.Storage) *SubscriptionController {
	return &SubscriptionController{stg: stg}
}

// Add new subscription to user with publication
func (scon *SubscriptionController) Add(user *logic.User, publication *logic.Publication) error {
	return nil
}

// Remove existing sybscription from user
func (scon *SubscriptionController) Remove(user *logic.User, publication *logic.Publication) error {
	return nil
}

// Update selected subscription
func (scon *SubscriptionController) Update(user *logic.User, publication *logic.Publication) error {
	return nil
}
