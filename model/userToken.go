package model

import (
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type UserToken struct {
	// todo  userId,accessToken, refreshToken,
	gorm.Model

	UserId string       `json:"user_id`
	Token  oauth2.Token `gorm:"foreignKey:UserId"`
}
