package models

import (
	"time"

	"gorm.io/gorm"
)

// DB MODELS
type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"`
	Name         string `json:"name" gorm:"type:varchar(50);not null"`
	UUID         string `json:"uuid" gorm:"type:uuid;unique;not null"`
	Login        string `json:"login" gorm:"type:varchar(50);unique;not null"`
	PasswordHash []byte `json:"-" gorm:"not null"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	EditedAt  time.Time      `json:"edited_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Question struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Text string `json:"text" gorm:"type:text;not null"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	EditedAt  time.Time      `json:"edited_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	Answers []Answer `json:"answers,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type Answer struct {
	ID         uint   `json:"id" gorm:"primaryKey"`
	QuestionID uint   `json:"question_id" gorm:"not null;index"` // индекс для быстрого поиска по вопросу
	UserID     string `json:"user_id" gorm:"type:uuid;not null"`
	User       User   `json:"user,omitempty" gorm:"foreignKey:UserID;references:UUID"`
	Text       string `json:"text" gorm:"type:text;not null"`

	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	EditedAt  time.Time      `json:"edited_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
