package controller

import (
	"sync"

	"github.com/aoyako/telegram_2ch_res_bot/storage"
)

// InfoController is an implementation of controller.Info
type InfoController struct {
	stg *storage.Storage
	m   sync.Mutex
}

// NewInfoController constructor of InfoController struct
func NewInfoController(stg *storage.Storage) *InfoController {
	return &InfoController{stg: stg}
}

// GetLastTimestamp returns time of the latest post
func (icon *InfoController) GetLastTimestamp() uint64 {
	return icon.stg.GetLastTimestamp()
}

// SetLastTimestamp sets time of the latest post
func (icon *InfoController) SetLastTimestamp(tsp uint64) {
	icon.m.Lock()
	last := icon.GetLastTimestamp()
	if last < tsp {
		icon.stg.SetLastTimestamp(tsp)
	}
	icon.m.Unlock()
}
