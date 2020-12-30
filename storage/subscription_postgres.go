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
	result := subsStorage.db.Create(publication)
	return result.Error
}

// Remove existing sybscription from user
func (subsStorage *SubscriptionPostgres) Remove(user *logic.User, publication *logic.Publication) error {
	result := subsStorage.db.Delete(publication)
	return result.Error
}

// Update selected subscription
func (subsStorage *SubscriptionPostgres) Update(user *logic.User, publication *logic.Publication) error {
	result := subsStorage.db.Save(publication)
	return result.Error
}

// GetSubsByUser returns list of user's subscriptions
func (subsStorage *SubscriptionPostgres) GetSubsByUser(user *logic.User) ([]logic.Publication, error) {
	var pubs []logic.Publication
	result := subsStorage.db.Model(&logic.Publication{}).Where("user_id = ?", user.ID).Find(&pubs)
	return pubs, result.Error
}

// GetAllSubs Returns all publications
func (subsStorage *SubscriptionPostgres) GetAllSubs() []logic.Publication {
	pubs := make([]logic.Publication, 0)
	subsStorage.db.Model(&logic.Publication{}).Find(&pubs)
	return pubs
}
