package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId               string
	FirstName            string
	LastName             string
	AccountType          string
	AllowedDataTransfer  uint64 //in Bytes
	ConsumedDataTransfer uint64 //in Bytes
	AllowedSpeed         int    // in MegaBytes
	AllowedThreads       int    // number of threads that can be assigned to a user at once
}
