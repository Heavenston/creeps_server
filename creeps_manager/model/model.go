package model

import (
	"time"

	"gorm.io/gorm"
)

type UserDiscordAuth struct {
	AccessToken  string
	TokenExpires time.Time
	RefreshToken string
	Scopes       string
}

type User struct {
	gorm.Model
	DiscordAuth UserDiscordAuth `gorm:"embedded;embeddedPrefix:discord_"`
}
