package repository

import (
	"github.com/google/uuid"
	"github.com/Est1ege/go-user-api/internal/domain/models"
)

// UserRepository определяет интерфейс для работы с хранилищем пользователей
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	GetAll() ([]*models.User, error)
}
