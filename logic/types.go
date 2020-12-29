package logic

import "gorm.io/gorm"

// User stores info about user
type User struct {
	gorm.Model
	ChatID    uint64        `gorm:"uniqueIndex"` // Telegram's chat id
	SubsCount uint          // Amount of current subscribtions
	Subs      []Publication `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // User's subscriptions
}

// Publication stores info about origin of data sent to user
type Publication struct {
	gorm.Model
	Board  string // 2ch board name
	Tags   string // Array of strings to search in thread title
	UserID uint   // Publication owner
	Type   string
}

// Info stores addition information about bot
type Info struct {
	gorm.Model
	LastPost uint64 // Time of the latest post
}
