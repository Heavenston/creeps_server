package model

import "gorm.io/gorm"

type UserModel struct {
    gorm.Model
    username string
}
