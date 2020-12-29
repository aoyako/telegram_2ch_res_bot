package controller

import "github.com/aoyako/telegram_2ch_res_bot/storage"

// InfoController is an implementation of controller.Info
type InfoController struct {
	stg *storage.Storage
}

// NewInfoController constructor of InfoController struct
func NewInfoController(stg *storage.Storage) *InfoController {
	return &InfoController{stg: stg}
}

// GetLastTimestamp returns time of the latest post
func (icon *InfoController) GetLastTimestamp() uint {
	return 0
}

// SetLastTimestamp sets time of the latest post
func (icon *InfoController) SetLastTimestamp(tsp uint) {

}
