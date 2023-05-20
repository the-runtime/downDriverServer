package model

import (
	"gorm.io/gorm"
	"time"
)

type SingleHistory struct {
	gorm.Model
	UserId     string    `json:"userid"`
	Filename   string    `json:"filename"`
	Filesize   uint64    `json:"size"`
	Startedat  time.Time `json:"startedat"`
	Finishedat time.Time `json:"finishedat"`
}

//type UserHistory struct {
//	gorm.Model
//	UserId      string
//	ListHistory []SingleHistory
//}
