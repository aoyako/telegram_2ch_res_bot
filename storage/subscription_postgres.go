package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// SubscriptionPostgres is an implementation of storage.Subscription
type SubscriptionPostgres struct {
	db *gorm.DB
}

// NewSubscriptionPostgres constructor of SubscriptionPostgres struct
func NewSubscriptionPostgres(db *gorm.DB) *SubscriptionPostgres {
	return &SubscriptionPostgres{
		db: db,
	}
}

// Add new subscription to user with publication
func (subsStorage *SubscriptionPostgres) Add(user *logic.User, publication *logic.Publication) error {
	return nil
}

// Remove existing sybscription from user
func (subsStorage *SubscriptionPostgres) Remove(user *logic.User, publication *logic.Publication) error {
	return nil
}

// Update selected subscription
func (subsStorage *SubscriptionPostgres) Update(user *logic.User, publication *logic.Publication) error {
	return nil
}
