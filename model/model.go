package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserId            string
	AllowedBandwidth  int
	AllowedSpeed      int
	ConsumedBandwidth int
}
