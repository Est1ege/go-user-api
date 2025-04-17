package service

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Est1ege/go-user-api/internal/domain/models"
)

// Создаем мок для UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_Create(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	
	input := models.CreateUserInput{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
	}
	
	// Case 1: Email already exists
	mockRepo.On("GetByEmail", input.Email).Return(&models.User{}, nil).Once()
	
	// Act
	user, err := service.Create(input)
	
	// Assert
	assert.Nil(t, user)
	assert.Equal(t, ErrEmailAlreadyExists, err)
	
	// Case 2: Successful creation
	mockRepo.On("GetByEmail", input.Email).Return(nil, errors.New("not found")).Once()
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil).Once()
	
	// Act
	user, err = service.Create(input)
	
	// Assert
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, input.Email, user.Email)
	assert.Equal(t, input.FirstName, user.FirstName)
	assert.Equal(t, input.LastName, user.LastName)
	
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetByID(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	
	id := uuid.New()
	expectedUser := &models.User{
		ID:        id,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	
	// Case 1: User found
	mockRepo.On("GetByID", id).Return(expectedUser, nil).Once()
	
	// Act
	user, err := service.GetByID(id)
	
	// Assert
	assert.Nil(t, err)
	assert.Equal(t, expectedUser, user)
	
	// Case 2: User not found
	mockRepo.On("GetByID", id).Return(nil, errors.New("user not found")).Once()
	
	// Act
	user, err = service.GetByID(id)
	
	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "user not found", err.Error())
	
	mockRepo.AssertExpectations(t)
}

func TestUserService_Update(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	
	id := uuid.New()
	existingUser := &models.User{
		ID:        id,
		Email:     "old@example.com",
		FirstName: "Old",
		LastName:  "Name",
	}
	
	input := models.UpdateUserInput{
		Email:     "new@example.com",
		FirstName: "New",
		LastName:  "Name",
	}
	
	// Case 1: User not found
	mockRepo.On("GetByID", id).Return(nil, errors.New("user not found")).Once()
	
	// Act
	user, err := service.Update(id, input)
	
	// Assert
	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.Equal(t, "user not found", err.Error())
	
	// Case 2: Email already exists
	mockRepo.On("GetByID", id).Return(existingUser, nil).Once()
	mockRepo.On("GetByEmail", input.Email).Return(&models.User{ID: uuid.New()}, nil).Once()
	
	// Act
	user, err = service.Update(id, input)
	
	// Assert
	assert.Nil(t, user)
	assert.Equal(t, ErrEmailAlreadyExists, err)
	
	// Case 3: Successful update
	mockRepo.On("GetByID", id).Return(existingUser, nil).Once()
	mockRepo.On("GetByEmail", input.Email).Return(nil, errors.New("not found")).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil).Once()
	
	// Act
	user, err = service.Update(id, input)
	
	// Assert
	assert.NotNil(t, user)
	assert.Nil(t, err)
	assert.Equal(t, input.Email, user.Email)
	assert.Equal(t, input.FirstName, user.FirstName)
	assert.Equal(t, input.LastName, user.LastName)
	
	mockRepo.AssertExpectations(t)
}

func TestUserService_Delete(t *testing.T) {
	// Arrange
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	
	id := uuid.New()
	
	// Case 1: Successful deletion
	mockRepo.On("Delete", id).Return(nil).Once()
	
	// Act
	err := service.Delete(id)
	
	// Assert
	assert.Nil(t, err)
	
	// Case 2: Error during deletion
	mockRepo.On("Delete", id).Return(errors.New("deletion error")).Once()
	
	// Act
	err = service.Delete(id)
	
	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, "deletion error", err.Error())
	
	mockRepo.AssertExpectations(t)
}