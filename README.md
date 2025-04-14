# ПВЗ API Service

Сервис для работы с ПВЗ (пунктами выдачи заказов).

## Функционал

- Аутентификация пользователей (JWT)
- Регистрация и вход
- Управление пунктами выдачи заказов (ПВЗ)
- Управление приёмкой товаров
- Добавление и удаление товаров в рамках приёмки

## Требования

- Docker и Docker Compose

## Запуск

1. Клонировать репозиторий:
```
git clone https://github.com/your-username/ops_service.git
cd ops_service
```

2. Запустить через Docker Compose:
```
docker-compose up --build
```

Сервис будет доступен по следующим адресам:
- HTTP API: http://localhost:8080
- gRPC API: localhost:3000
- Prometheus метрики: http://localhost:9000/metrics

## API Endpoints

### Аутентификация

- `POST /register` - Регистрация нового пользователя
- `POST /login` - Вход в систему
- `GET /dummyLogin?role=client|moderator` - Вход для тестирования (без пароля)

### ПВЗ

- `POST /pickup-points` - Создание нового ПВЗ (только для модератора)
- `GET /pickup-points` - Получение списка ПВЗ с пагинацией

### Приёмка товаров

- `POST /receipts` - Создание новой приёмки товаров
- `POST /receipts/:id/close` - Закрытие приёмки
- `POST /receipts/products` - Добавление товара в рамках приёмки
- `DELETE /receipts/products` - Удаление последнего добавленного товара

## Архитектура

Проект организован в соответствии с принципами чистой архитектуры:

- `cmd/app` - Точка входа приложения
- `internal/api` - API обработчики
- `internal/auth` - Аутентификация и авторизация
- `internal/config` - Конфигурация приложения
- `internal/db` - Подключение к базе данных
- `internal/middleware` - Промежуточные обработчики
- `internal/model` - Модели данных
- `internal/repository` - Репозитории для работы с данными
- `internal/service` - Бизнес-логика
- `migrations` - SQL миграции

## Примеры запросов

### Регистрация пользователя

```
POST /register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "role": "client"
}
```

### Получение токена (dummy login)

```
GET /dummyLogin?role=moderator
```

### Создание ПВЗ (требуется токен модератора)

```
POST /pickup-points
Authorization: Bearer <token>
Content-Type: application/json

{
  "city": "Москва"
}
``` 
