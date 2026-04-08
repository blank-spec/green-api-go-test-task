# Go backend for Green API test task

Сервис поднимает REST API поверх методов Green API, проксирует запросы в Green API и отдает одну HTML-страницу с формами для работы через браузер.

## Что реализовано

- `GET /healthz` - healthcheck
- `GET /` - HTML-страница, где пользователь сам вводит `idInstance` и `apiTokenInstance`
- `POST /api/v1/settings` - получить настройки инстанса
- `POST /api/v1/state` - получить состояние инстанса
- `POST /api/v1/messages/text` - отправить текстовое сообщение
- `POST /api/v1/messages/file` - отправить файл по URL

## Переменные окружения

- `HTTP_ADDR` - адрес HTTP-сервера, по умолчанию `:8080`
- `GREEN_API_BASE_URL` - базовый URL Green API, по умолчанию `https://api.green-api.com`
- `GREEN_API_REQUEST_TIMEOUT` - timeout для запросов к upstream, по умолчанию `15s`

## Запуск

```powershell
$env:GOTELEMETRY='off'
go run .\\cmd\\server
```

После запуска открой `http://localhost:8080/` и введи `idInstance` и `apiTokenInstance` прямо на странице.

## Docker

### Сборка образа

```powershell
docker build -t green-api-form .
```

### Запуск контейнера

```powershell
docker run --rm -p 8080:8080 green-api-form
```

### Запуск через Docker Compose

```powershell
docker compose up -d --build
```
