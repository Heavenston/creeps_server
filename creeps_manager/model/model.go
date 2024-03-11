package model

import (
	"time"

	"github.com/Heavenston/creeps_server/creeps_manager/model/discordmodel"
	"gorm.io/gorm"
)

type UserDiscordAuth struct {
	AccessToken  string
	TokenExpires time.Time
	RefreshToken string
	Scope        string
}

type User struct {
	gorm.Model
	DiscordId   string            `gorm:"uniqueIndex"`
	DiscordAuth UserDiscordAuth   `gorm:"embedded;embeddedPrefix:discord_"`
	DiscordUser discordmodel.User `gorm:"serializer:json"`

	RoleID *int
	Role   *Role `gorm:"constraint:OnDelete:SET NULL;"`
}

type Role struct {
	gorm.Model
	Name string
}

type GameConfig struct {
	CanJoinAfterStart bool
	Private           bool
	IsLocal           bool
}

type Game struct {
	gorm.Model
	Name string

	CreatorID int
	Creator   *User

	Players []User `gorm:"many2many:game_players;"`

	Config GameConfig `gorm:"embedded"`

	StartedAt *time.Time
	EndedAt   *time.Time
}
