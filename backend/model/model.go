package model

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;unique;not null"`
	Password  string    `gorm:"size:255;not null"`
	Bio       *string   `gorm:"size:255"`
	Image     *string   `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Post struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	Content   string    `gorm:"size:400;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Follow struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	UserID     uint      `gorm:"not null"`
	FollowerID uint      `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
