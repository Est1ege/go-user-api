package service

import (
	"errors"
	"github.com/google/uuid"
	"github.com/Est1ege/go-user-api/internal/domain/models"
	"github.com/Est1ege/go-user-api/internal/repository/postgres"
	"golang.org/x/crypto/bcrypt"
)

// UserService представляет сервис для работы с пользователями
type UserService struct {
	userRepo *postgres.UserRepository
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo *postgres.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Create создает нового пользователя
func (s *UserService) Create(input models.CreateUserInput) (*models.User, error) {
	// Проверяем, существует ли пользователь с таким email
	existingUser, err := s.userRepo.GetByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &models.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Password:  string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GetByID получает пользователя по ID
func (s *UserService) GetByID(id uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

// Update обновляет данные пользователя
func (s *UserService) Update(id uuid.UUID, input models.UpdateUserInput) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Обновляем поля, если они были предоставлены
	if input.Email != "" && input.Email != user.Email {
		// Проверяем, не занят ли новый email
		existingUser, err := s.userRepo.GetByEmail(input.Email)
		if err == nil && existingUser != nil && existingUser.ID != id {
			return nil, ErrEmailAlreadyExists
		}
		user.Email = input.Email
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}

	if input.LastName != "" {
		user.LastName = input.LastName
	}

	if input.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Delete удаляет пользователя
func (s *UserService) Delete(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}

// GetAll получает список всех пользователей
func (s *UserService) GetAll() ([]*models.User, error) {
	return s.userRepo.GetAll()
}

// Определение ошибок
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
)