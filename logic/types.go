package logic

import "gorm.io/gorm"

// User stores info about user
type User struct {
	gorm.Model
	ChatID    uint // Telegram's chat id
	SubsCount uint // Amount of current subscribtions
}

// Publication stores info about origin of data sent to user
type Publication struct {
	gorm.Model
	Board string   // 2ch board
	Tags  []string // Array of strings to search in thread title
}

// Subscription represents connection between user and publication
type Subscription struct {
	gorm.Model
	UserID        uint
	PublicationID uint
}

// Info stores addition information about bot
type Info struct {
	gorm.Model
	LastPost uint // Time of the latest post
}
