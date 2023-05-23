package modelTask

import (
	"gorm.io/gorm"
	modelUser "mbf5923.com/todo/domain/user/models"
	"time"
)

type EntityTask struct {
	ID          uint                  `gorm:"primaryKey;"`
	UserId      uint                  `gorm:"type:int;not null"`
	User        modelUser.EntityUsers `gorm:"foreignKey:UserId"`
	Title       string                `gorm:"type:varchar(255);not null"`
	Description string                `gorm:"type:varchar(500);not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (entity *EntityTask) BeforeCreate(db *gorm.DB) error {
	entity.CreatedAt = time.Now().Local()
	return nil
}

func (entity *EntityTask) BeforeUpdate(db *gorm.DB) error {
	entity.UpdatedAt = time.Now().Local()
	return nil
}
