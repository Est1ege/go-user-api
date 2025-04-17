package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User представляет модель пользователя
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Email     string    `gorm:"type:varchar(100);unique_index" json:"email" binding:"required,email"`
	FirstName string    `gorm:"type:varchar(100)" json:"first_name" binding:"required"`
	LastName  string    `gorm:"type:varchar(100)" json:"last_name" binding:"required"`
	Password  string    `gorm:"type:varchar(255)" json:"-"` // Не отправляем пароль в JSON
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateUserInput определяет структуру для создания пользователя
type CreateUserInput struct {
	Email     string `json:"email" form:"email" binding:"required,email"`
	FirstName string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`
	Password  string `json:"password" form:"password" binding:"required,min=8"`
}

// UpdateUserInput определяет структуру для обновления данных пользователя
type UpdateUserInput struct {
    Email     string `json:"email" form:"email" binding:"omitempty,email"`
    FirstName string `json:"first_name" form:"first_name"`
    LastName  string `json:"last_name" form:"last_name"`
    Password  string `json:"password" form:"password" binding:"omitempty,min=8"`
}

// BeforeCreate - хук GORM, который выполняется перед созданием записи
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}