package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/Est1ege/go-user-api/internal/domain/models"
	"github.com/Est1ege/go-user-api/internal/service"
)

// MockUserService имитирует сервис пользователя для тестирования
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll() ([]*models.User, error) {
    args := m.Called()
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]*models.User), args.Error(1)
}

// Убедимся что MockUserService реализует service.UserServiceInterface
var _ service.UserServiceInterface = (*MockUserService)(nil)


func (m *MockUserService) Create(input models.CreateUserInput) (*models.User, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) GetByID(id uuid.UUID)  (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Update(id uuid.UUID, input models.UpdateUserInput) (*models.User, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserService) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() (*gin.Engine, *MockUserService) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)
	
	// Настраиваем маршруты
	userRoutes := router.Group("/users")
	{
		userRoutes.POST("", handler.Create)
		userRoutes.GET("/:id", handler.GetByID)
		userRoutes.PUT("/:id", handler.Update)
		userRoutes.DELETE("/:id", handler.Delete)
	}
	
	return router, mockService
}

func TestUserHandler_Create(t *testing.T) {
	// Arrange
	router, mockService := setupTestRouter()
	
	// Test case: успешное создание пользователя
	input := models.CreateUserInput{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password123",
	}
	
	createdUser := &models.User{
		ID:        uuid.New(),
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	
	mockService.On("Create", mock.AnythingOfType("models.CreateUserInput")).Return(createdUser, nil).Once()
	
	// Преобразуем входные данные в JSON
	jsonInput, _ := json.Marshal(input)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, createdUser.ID, response.ID)
	assert.Equal(t, createdUser.Email, response.Email)
	
	// Test case: ошибка валидации
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(`{"email": "invalid"}`)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	// Test case: ошибка "email уже существует"
	mockService.On("Create", mock.AnythingOfType("models.CreateUserInput")).Return(nil, service.ErrEmailAlreadyExists).Once()
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusConflict, w.Code)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_GetByID(t *testing.T) {
	// Arrange
	router, mockService := setupTestRouter()
	
	id := uuid.New()
	user := &models.User{
		ID:        id,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}
	
	// Test case: успешное получение пользователя
	mockService.On("GetByID", id).Return(user, nil).Once()
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+id.String(), nil)
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Email, response.Email)
	
	// Test case: пользователь не найден
	mockService.On("GetByID", id).Return(nil, errors.New("user not found")).Once()
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/"+id.String(), nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	// Test case: некорректный ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/users/invalid-id", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_Update(t *testing.T) {
	// Arrange
	router, mockService := setupTestRouter()
	
	id := uuid.New()
	input := models.UpdateUserInput{
		Email:     "updated@example.com",
		FirstName: "Updated",
		LastName:  "Name",
	}
	
	updatedUser := &models.User{
		ID:        id,
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}
	
	// Test case: успешное обновление пользователя
	mockService.On("Update", id, mock.AnythingOfType("models.UpdateUserInput")).Return(updatedUser, nil).Once()
	
	// Преобразуем входные данные в JSON
	jsonInput, _ := json.Marshal(input)
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/"+id.String(), bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, updatedUser.ID, response.ID)
	assert.Equal(t, updatedUser.Email, response.Email)
	
	// Test case: ошибка "email уже существует"
	mockService.On("Update", id, mock.AnythingOfType("models.UpdateUserInput")).Return(nil, service.ErrEmailAlreadyExists).Once()
	
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/users/"+id.String(), bytes.NewBuffer(jsonInput))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusConflict, w.Code)
	
	mockService.AssertExpectations(t)
}

func TestUserHandler_Delete(t *testing.T) {
	// Arrange
	router, mockService := setupTestRouter()
	
	id := uuid.New()
	
	// Test case: успешное удаление пользователя
	mockService.On("Delete", id).Return(nil).Once()
	
	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/users/"+id.String(), nil)
	router.ServeHTTP(w, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Test case: некорректный ID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/users/invalid-id", nil)
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	mockService.AssertExpectations(t)
}