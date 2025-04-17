package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Est1ege/go-user-api/internal/domain/models"
	"github.com/Est1ege/go-user-api/internal/service"
)

// WebHandler представляет обработчики для веб-интерфейса
type WebHandler struct {
	userService *service.UserService
}

// NewWebHandler создает новый экземпляр WebHandler
func NewWebHandler(userService *service.UserService) *WebHandler {
	return &WebHandler{userService: userService}
}

// Index отображает страницу со списком пользователей
func (h *WebHandler) Index(c *gin.Context) {
	// Получение всех пользователей
	users, err := h.userService.GetAll()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "users/index.html", gin.H{
			"Error": "Ошибка при получении списка пользователей: " + err.Error(),
		})
		return
	}

	// Получение флеш-сообщений из сессии
	success, _ := c.Get("success")

	c.HTML(http.StatusOK, "users/index.html", gin.H{
		"Users":   users,
		"Success": success,
	})
}

// Create создает нового пользователя через веб-форму
func (h *WebHandler) Create(c *gin.Context) {
	var input models.CreateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": "Ошибка валидации данных: " + err.Error(),
		})
		return
	}

	_, err := h.userService.Create(input)
	if err != nil {
		errorMessage := "Ошибка при создании пользователя"
		if err == service.ErrEmailAlreadyExists {
			errorMessage = "Email уже используется"
		}
		
		// Получаем всех пользователей для отображения на странице
		users, _ := h.userService.GetAll()
		
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": errorMessage,
			"Users": users,
		})
		return
	}

	// Устанавливаем флеш-сообщение
	c.Set("success", "Пользователь успешно создан")
	c.Redirect(http.StatusSeeOther, "/web/users")
}

// Update обновляет данные пользователя через веб-форму
func (h *WebHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": "Некорректный ID пользователя",
		})
		return
	}

	var input models.UpdateUserInput
	if err := c.ShouldBind(&input); err != nil {
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": "Ошибка валидации данных: " + err.Error(),
		})
		return
	}

	_, err = h.userService.Update(id, input)
	if err != nil {
		errorMessage := "Ошибка при обновлении пользователя"
		if err == service.ErrEmailAlreadyExists {
			errorMessage = "Email уже используется"
		}
		
		// Получаем всех пользователей для отображения на странице
		users, _ := h.userService.GetAll()
		
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": errorMessage,
			"Users": users,
		})
		return
	}

	// Устанавливаем флеш-сообщение
	c.Set("success", "Пользователь успешно обновлен")
	c.Redirect(http.StatusSeeOther, "/web/users")
}

// Delete удаляет пользователя через веб-форму
func (h *WebHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": "Некорректный ID пользователя",
		})
		return
	}

	if err := h.userService.Delete(id); err != nil {
		// Получаем всех пользователей для отображения на странице
		users, _ := h.userService.GetAll()
		
		c.HTML(http.StatusBadRequest, "users/index.html", gin.H{
			"Error": "Ошибка при удалении пользователя: " + err.Error(),
			"Users": users,
		})
		return
	}

	// Устанавливаем флеш-сообщение
	c.Set("success", "Пользователь успешно удален")
	c.Redirect(http.StatusSeeOther, "/web/users")
}