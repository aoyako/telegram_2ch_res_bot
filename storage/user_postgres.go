package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// UserPostgres is an implementation of storage.User
type UserPostgres struct {
	db  *gorm.DB
	cfg *InitDatabase
}

// NewUserPostgres constructor of UserPostgres struct
func NewUserPostgres(db *gorm.DB, cfg *InitDatabase) *UserPostgres {
	return &UserPostgres{
		db:  db,
		cfg: cfg,
	}
}

// Register adds user in databse
func (userStorage *UserPostgres) Register(user *logic.User) error {
	var tUser logic.User
	exists := userStorage.db.Where("chat_id = ?", user.ChatID).First(&tUser)
	if exists.Error != nil {
		result := userStorage.db.Create(user)

		if result.Error != nil {
			return result.Error
		}

		// Adds admin record, if admin added
		if _, contains := contains(userStorage.cfg.Admin, user.ChatID); contains {
			result := userStorage.db.Create(&logic.Admin{
				UserID: uint64(user.ID),
			})

			if result.Error != nil {
				return result.Error
			}
		}
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

// GetUserByID returns user by it's id
func (userStorage *UserPostgres) GetUserByID(userID uint) (*logic.User, error) {
	var user logic.User
	result := userStorage.db.Where("id = ?", userID).First(&user)

	return &user, result.Error
}

// GetUsersByPublication returns subscribers of publication
func (userStorage *UserPostgres) GetUsersByPublication(pub *logic.Publication) ([]logic.User, error) {
	var users []logic.User
	result := userStorage.db.Model(&pub).Association("Users").Find(&users)

	return users, result
}

// IsUserAdmin checks if user has administrator privileges
func (userStorage *UserPostgres) IsUserAdmin(user *logic.User) bool {
	var count int64
	userStorage.db.Find(&logic.Admin{}).Where("user_id = ?", user.ID).Count(&count)
	return count != 0
}

// IsChatAdmin checks if user has administrator privileges by chatID
func (userStorage *UserPostgres) IsChatAdmin(chatID uint) bool {
	var user logic.User
	userStorage.db.Find(&logic.User{}).Where("chat_id = ?", chatID).First(&user)
	return userStorage.IsUserAdmin(&user)
}

func contains(slice []uint64, val uint64) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}
