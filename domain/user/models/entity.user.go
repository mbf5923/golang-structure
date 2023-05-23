package modelUser

import (
	"gorm.io/gorm"
	util "mbf5923.com/todo/utils"
	"time"
)

type EntityUsers struct {
	ID        uint   `gorm:"primaryKey;"`
	Fullname  string `gorm:"type:varchar(255);unique;not null"`
	Email     string `gorm:"type:varchar(255);unique;not null"`
	Password  string `gorm:"type:varchar(255);not null"`
	Active    bool   `gorm:"type:bool;default:false"`
	ApiKey    string `gorm:"type:varchar(255);unique;null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (entity *EntityUsers) BeforeCreate(db *gorm.DB) error {
	entity.Password = util.HashPassword(entity.Password)
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityUsers) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
