package handlers

import (
    "log"  // Добавьте импорт для логирования
    "net/http"
    
    "github.com/gin-contrib/sessions"  // Добавьте импорт для сессий
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
    // Получение сессии
    session := sessions.Default(c)
    
    // Получаем флеш-сообщения
    success := session.Flashes("success")
    error := session.Flashes("error")
    
    // Сохраняем сессию
    session.Save()
    
    // Получение всех пользователей
    users, err := h.userService.GetAll()
    if err != nil {
        log.Printf("Ошибка при получении списка пользователей: %v", err)
        c.HTML(http.StatusInternalServerError, "index.html", gin.H{
            "Error": "Ошибка при получении списка пользователей: " + err.Error(),
        })
        return
    }
    
    // Формируем данные для шаблона
    data := gin.H{"Users": users}
    
    if len(success) > 0 {
        data["Success"] = success[0]
    }
    
    if len(error) > 0 {
        data["Error"] = error[0]
    }
    
    // Рендерим шаблон
    c.HTML(http.StatusOK, "index.html", data)
}

// Create создает нового пользователя через веб-форму
func (h *WebHandler) Create(c *gin.Context) {
    log.Printf("Создание пользователя, метод: %s", c.Request.Method)
    log.Printf("Данные формы: %+v", c.Request.Form)
    
    // Разбор формы
    if err := c.Request.ParseForm(); err != nil {
        log.Printf("Ошибка при разборе формы: %v", err)
    }
    log.Printf("Форма после разбора: %+v", c.Request.Form)
    
    var input models.CreateUserInput
    if err := c.ShouldBind(&input); err != nil {
        log.Printf("Ошибка привязки данных: %v", err)
        
        // Получаем всех пользователей для отображения на странице
        users, _ := h.userService.GetAll()
        
        c.HTML(http.StatusOK, "index.html", gin.H{
            "Error": "Ошибка валидации данных: " + err.Error(),
            "Users": users,
        })
        return
    }
    
    log.Printf("Данные для создания пользователя: %+v", input)
    
    _, err := h.userService.Create(input)
    if err != nil {
        log.Printf("Ошибка при создании пользователя: %v", err)
        
        errorMessage := "Ошибка при создании пользователя"
        if err == service.ErrEmailAlreadyExists {
            errorMessage = "Email уже используется"
        }
        
        // Получаем всех пользователей для отображения на странице
        users, _ := h.userService.GetAll()
        
        c.HTML(http.StatusOK, "index.html", gin.H{
            "Error": errorMessage,
            "Users": users,
        })
        return
    }
    
    // Устанавливаем флеш-сообщение в сессии
    session := sessions.Default(c)
    session.AddFlash("Пользователь успешно создан", "success")
    session.Save()
    
    c.Redirect(http.StatusSeeOther, "/web/users")
}

// Update обновляет данные пользователя через веб-форму
func (h *WebHandler) Update(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        log.Printf("Некорректный ID пользователя: %v", err)
        
        // Добавляем сообщение об ошибке в сессию
        session := sessions.Default(c)
        session.AddFlash("Некорректный ID пользователя", "error")
        session.Save()
        
        c.Redirect(http.StatusSeeOther, "/web/users")
        return
    }
    
    var input models.UpdateUserInput
    if err := c.ShouldBind(&input); err != nil {
        log.Printf("Ошибка валидации данных при обновлении: %v", err)
        
        // Добавляем сообщение об ошибке в сессию
        session := sessions.Default(c)
        session.AddFlash("Ошибка валидации данных: "+err.Error(), "error")
        session.Save()
        
        c.Redirect(http.StatusSeeOther, "/web/users")
        return
    }
    
    _, err = h.userService.Update(id, input)
    if err != nil {
        log.Printf("Ошибка при обновлении пользователя: %v", err)
        
        errorMessage := "Ошибка при обновлении пользователя"
        if err == service.ErrEmailAlreadyExists {
            errorMessage = "Email уже используется"
        }
        
        // Добавляем сообщение об ошибке в сессию
        session := sessions.Default(c)
        session.AddFlash(errorMessage, "error")
        session.Save()
        
        c.Redirect(http.StatusSeeOther, "/web/users")
        return
    }
    
    // Устанавливаем флеш-сообщение в сессии
    session := sessions.Default(c)
    session.AddFlash("Пользователь успешно обновлен", "success")
    session.Save()
    
    c.Redirect(http.StatusSeeOther, "/web/users")
}

// Delete удаляет пользователя через веб-форму
func (h *WebHandler) Delete(c *gin.Context) {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        log.Printf("Некорректный ID пользователя при удалении: %v", err)
        
        // Добавляем сообщение об ошибке в сессию
        session := sessions.Default(c)
        session.AddFlash("Некорректный ID пользователя", "error")
        session.Save()
        
        c.Redirect(http.StatusSeeOther, "/web/users")
        return
    }
    
    if err := h.userService.Delete(id); err != nil {
        log.Printf("Ошибка при удалении пользователя: %v", err)
        
        // Добавляем сообщение об ошибке в сессию
        session := sessions.Default(c)
        session.AddFlash("Ошибка при удалении пользователя: "+err.Error(), "error")
        session.Save()
        
        c.Redirect(http.StatusSeeOther, "/web/users")
        return
    }
    
    // Устанавливаем флеш-сообщение в сессии
    session := sessions.Default(c)
    session.AddFlash("Пользователь успешно удален", "success")
    session.Save()
    
    c.Redirect(http.StatusSeeOther, "/web/users")
}