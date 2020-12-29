package storage

import (
	"github.com/aoyako/telegram_2ch_res_bot/logic"
	"gorm.io/gorm"
)

// InfoPostgres is an implementation of storage.Info
type InfoPostgres struct {
	db *gorm.DB
}

// NewInfoPostgres constructor of InfoPostgres struct
func NewInfoPostgres(db *gorm.DB) *InfoPostgres {
	return &InfoPostgres{
		db: db,
	}
}

// GetLastTimestamp returns time of the latest post
func (infoStorage *InfoPostgres) GetLastTimestamp() uint64 {
	var info logic.Info
	infoStorage.db.First(&info)

	return info.LastPost
}

// SetLastTimestamp sets time of the latest post
func (infoStorage *InfoPostgres) SetLastTimestamp(timestamp uint64) {
	var info logic.Info
	infoStorage.db.First(&info)
	info.LastPost = timestamp

	infoStorage.db.Save(&info)
}
