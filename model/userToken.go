package model

import (
	"gorm.io/gorm"
	"time"
)

type UserToken struct {
	// todo  userId,accessToken, refreshToken,
	gorm.Model

	UserId       string
	AuthCode     string
	AccessToken  string
	TokenType    string
	RefreshToken string
	Expiry       time.Time
	//Token  oauth2.Token `gorm:"foreignKey:UserId"`
}
