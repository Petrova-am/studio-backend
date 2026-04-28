# Студия растяжки - Backend
Микросервисная архитектура для студии растяжки на Go + PostgreSQL.

##  Структура проекта
studio-backend/
├── cmd/
│ ├── auth-service/ # Сервис авторизации (порт 8081)
│ └── booking-service/ # Сервис бронирования (порт 8082)
├── internal/
│ ├── auth/ # Логика авторизации
│ ├── booking/ # Логика бронирования
│ └── models/ # Модели данных
├── pkg/
│ ├── database/ # Подключение к БД
│ ├── logger/ # Логирование (zerolog)
│ └── middleware/ # CORS, JWT
├── migrations/ # Миграции БД (3 файла)
└── docs/ # Swagger документация

##  Запуск проекта

### 1. Установка зависимостей
```bash
go mod tidy
```

### 2. Запуск PostgreSQL (Docker)
```
docker run --name postgres -e POSTGRES_PASSWORD=123987 -e POSTGRES_DB=studio_db -p 5432:5432 -d postgres:15
```

### 3. Миграции
```bash
migrate -path ./migrations -database "postgres://postgres:123987@localhost:5432/studio_db?sslmode=disable" up
```

### 4. Запуск сервисов

Auth сервис:
```bash
go run cmd/auth-service/main.go
```
Booking сервис:
```bash
go run cmd/booking-service/main.go
```
## Тестирование

### Все тесты с покрытием
```bash
go test -cover ./...
```

### Auth сервис
```bash
go test -cover ./cmd/auth-service/...
```

### Тесты с моками
```bash
go test -cover ./internal/auth/...
```

##  API Endpoints

### Auth Service (порт 8081)

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/auth/login` | Вход в систему |
| POST | `/api/auth/register` | Регистрация |
| GET | `/api/auth/validate` | Проверка токена |

### Booking Service (порт 8082)

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/bookings` | Создать бронирование |
| GET | `/api/bookings/user/{id}` | Бронирования пользователя |
| PUT | `/api/bookings/{id}/status` | Обновить статус |
| GET | `/api/bookings/available-slots` | Свободные места |

##  Миграции

| Файл | Описание |
|------|----------|
| `001_create_users_table.up.sql` | Таблица пользователей |
| `002_create_trainers_table.up.sql` | Таблица тренеров |
| `003_create_bookings_table.up.sql` | Таблица бронирований |

##  Технологии

| Технология | Версия | Назначение |
|------------|--------|------------|
| **Go** | 1.21 | Язык программирования |
| **PostgreSQL** | 15 | База данных |
| **JWT** | v5 | Авторизация и аутентификация |
| **Zerolog** | v1.32 | Логирование |
| **Testify** | v1.8 | Тесты и моки |
| **golang-migrate** | v4 | Управление миграциями |
| **Swagger** | v1.16 | Документация API |

##  Импорт модулей

```go
import "github.com/Petrova-am/studio-backend/pkg/database"
import "github.com/Petrova-am/studio-backend/pkg/logger"
import "github.com/Petrova-am/studio-backend/pkg/middleware"

Автор
Петрова Александра

## Сохранить и отправить:
```cmd
git add README.md
git commit -m "Add full README"
git push