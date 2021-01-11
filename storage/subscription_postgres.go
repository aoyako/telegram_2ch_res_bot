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
	subsStorage.db.Model(publication).Association("Users").Append(user)
	return result.Error
}

// AddDefault creates default publication
func (subsStorage *SubscriptionPostgres) AddDefault(publication *logic.Publication) error {
	publication.IsDefault = true
	result := subsStorage.db.Create(publication)
	return result.Error
}

// Remove existing sybscription from user
func (subsStorage *SubscriptionPostgres) Remove(publication *logic.Publication) error {
	result := subsStorage.db.Delete(publication)
	return result.Error
}

// Connect (subscribe) user to publication
func (subsStorage *SubscriptionPostgres) Connect(user *logic.User, publication *logic.Publication) error {
	result := subsStorage.db.Model(publication).Association("Users").Append(user)
	return result
}

// Disonnect (unsubscribe) user from publication
func (subsStorage *SubscriptionPostgres) Disonnect(user *logic.User, publication *logic.Publication) error {
	result := subsStorage.db.Model(publication).Association("Users").Delete(user)
	return result
}

// Update selected subscription
func (subsStorage *SubscriptionPostgres) Update(user *logic.User, publication *logic.Publication) error {
	result := subsStorage.db.Save(publication)
	return result.Error
}

// GetSubsByUser returns list of user's subscriptions
func (subsStorage *SubscriptionPostgres) GetSubsByUser(user *logic.User) ([]logic.Publication, error) {
	var pubs []logic.Publication
	result := subsStorage.db.Model(user).Association("Subs").Find(&pubs)
	return pubs, result
}

// GetAllSubs returns all publications
func (subsStorage *SubscriptionPostgres) GetAllSubs() []logic.Publication {
	pubs := make([]logic.Publication, 0)
	subsStorage.db.Model(&logic.Publication{}).Find(&pubs)
	return pubs
}

// GetAllDefaultSubs returns all default publications
func (subsStorage *SubscriptionPostgres) GetAllDefaultSubs() []logic.Publication {
	var pubs []logic.Publication
	subsStorage.db.Model(&logic.Publication{}).Where("is_default = ?", true).Find(&pubs)
	return pubs
}
