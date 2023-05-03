package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId               string
	AllowedDataTransfer  uint64 //in Bytes
	ConsumedDataTransfer uint64 //in Bytes
	AllowedSpeed         int    // in MegaBytes
}
