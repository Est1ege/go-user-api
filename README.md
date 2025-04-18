# Go User API

REST API на языке Go для управления пользователями с использованием PostgreSQL.

## Возможности

- Создание, получение, обновление и удаление пользователей
- Валидация входящих данных
- Хеширование паролей
- Модульная архитектура
- Контейнеризация с помощью Docker и Docker Compose
- Тесты для бизнес-логики и API

## Структура проекта

Проект построен с использованием чистой архитектуры:

- `cmd/api` - точка входа в приложение
- `internal` - внутренний код приложения
  - `api` - обработчики HTTP и маршрутизация
  - `config` - конфигурация приложения
  - `domain` - модели предметной области
  - `repository` - слой для работы с базой данных
  - `service` - бизнес-логика
- `pkg` - многоразовые пакеты, которые могут быть использованы другими приложениями

## Технологии

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin) - HTTP фреймворк
- [GORM](https://gorm.io/) - ORM для Go
- [PostgreSQL](https://www.postgresql.org/) - реляционная база данных
- [Docker](https://www.docker.com/) - для контейнеризации
- [Testify](https://github.com/stretchr/testify) - инструменты для тестирования

## API Endpoints

| Метод | Endpoint | Описание |
| --- | --- | --- |
| POST | /api/v1/users | Создание нового пользователя |
| GET | /api/v1/users/:id | Получение информации о пользователе по ID |
| PUT | /api/v1/users/:id | Обновление данных пользователя |
| DELETE | /api/v1/users/:id | Удаление пользователя |

## Локальный запуск

### Предварительные требования

- Go 1.18 или выше
- Docker и Docker Compose

### Запуск с помощью Docker Compose

```bash
# Клонирование репозитория
git clone https://github.com/yourusername/go-user-api.git
cd go-user-api

# Запуск приложения и базы данных с помощью Docker Compose
docker-compose up -d
```

API будет доступен по адресу http://localhost:8080

### Локальный запуск для разработки

```bash
# Настройка переменных окружения
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=user_api

# Запуск PostgreSQL в Docker
docker run -d -p 5432:5432 --name postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=user_api \
  postgres:15-alpine

# Запуск приложения
go run cmd/api/main.go
```

## Запуск тестов

```bash
# Запуск всех тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Генерация отчета о покрытии
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Пример использования API

### Создание пользователя

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "password": "securepassword"
  }'
```

### Получение информации о пользователе

```bash
curl -X GET http://localhost:8080/api/v1/users/YOUR_USER_ID
```

### Обновление данных пользователя

```bash
curl -X PUT http://localhost:8080/api/v1/users/YOUR_USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John Updated",
    "last_name": "Doe Updated"
  }'
```

### Удаление пользователя

```bash
curl -X DELETE http://localhost:8080/api/v1/users/YOUR_USER_ID
```

## Архитектура проекта

Проект построен с использованием подхода чистой архитектуры и принципа инверсии зависимостей:

1. **Доменный слой**: Содержит модели и бизнес-логику предметной области.
   - `models` - определяет сущности и структуры данных

2. **Слой репозитория**: Отвечает за доступ к данным.
   - `repository` - определяет интерфейсы для работы с данными
   - `repository/postgres` - конкретная реализация для PostgreSQL

3. **Сервисный слой**: Содержит бизнес-логику приложения.
   - `service` - реализует бизнес-правила и координирует работу с репозиториями

4. **Слой API**: Взаимодействует с клиентами.
   - `api/handlers` - обработчики HTTP-запросов
   - `api/middleware` - промежуточное ПО
   - `api/routes` - определение маршрутов

Взаимодействие между слоями организовано через интерфейсы, что обеспечивает:
- Слабую связанность компонентов
- Легкую тестируемость с использованием моков
- Возможность замены реализаций без изменения интерфейсов

## Лицензия

MIT