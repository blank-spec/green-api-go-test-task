# Green API Fiber Backend

Небольшой Go/Fiber backend для работы с [Green API](https://green-api.com/) через одну встроенную HTML-страницу.

Сервис поднимает UI на `/`, принимает `idInstance` и `apiTokenInstance` прямо из формы и проксирует запросы в Green API без хранения credentials на сервере.

## Features

- Одна HTML-страница для ручного тестирования
- `getSettings`
- `getStateInstance`
- `sendMessage`
- `sendFileByUrl`
- Валидация входных данных на backend
- Docker и Docker Compose для быстрого деплоя
- Fiber как HTTP transport layer

## Stack

- Go `1.25`
- Fiber `v2`
- Embedded HTML UI
- Docker

## Project Structure

```text
cmd/server                 entrypoint
internal/config            runtime config
internal/greenapi          Green API client
internal/httpapi           Fiber handlers, validation, embedded UI
```

## API

### UI

- `GET /` — встроенная HTML-страница
- `GET /healthz` — healthcheck

### Backend endpoints

- `POST /api/v1/settings`
- `POST /api/v1/state`
- `POST /api/v1/messages/text`
- `POST /api/v1/messages/file`

Все методы, кроме `GET /` и `GET /healthz`, требуют:

```json
{
  "idInstance": "1101000001",
  "apiTokenInstance": "your_token"
}
```

Если credentials не переданы, backend не выполнит запрос.

## Run Locally

```powershell
$env:GOTELEMETRY='off'
go run .\cmd\server
```

После запуска открой [http://localhost:8080/](http://localhost:8080/).

## Environment Variables

- `HTTP_ADDR` — адрес HTTP-сервера, default: `:8080`
- `GREEN_API_BASE_URL` — базовый URL Green API, default: `https://api.green-api.com`
- `GREEN_API_REQUEST_TIMEOUT` — timeout запросов к upstream, default: `15s`

## Docker

### Build

```powershell
docker build -t green-api-fiber-backend .
```

### Run

```powershell
docker run --rm -p 8080:8080 green-api-fiber-backend
```

### Run With Custom Base URL

```powershell
docker run --rm -p 8080:8080 `
  -e GREEN_API_BASE_URL=https://api.green-api.com `
  green-api-fiber-backend
```

### Docker Compose

```powershell
docker compose up -d --build
```

## Example Requests

### Get Settings

```powershell
$body = @{
    idInstance = '1101000001'
    apiTokenInstance = 'your_token'
} | ConvertTo-Json

Invoke-RestMethod `
    -Method Post `
    -Uri 'http://localhost:8080/api/v1/settings' `
    -ContentType 'application/json' `
    -Body $body
```

### Get State

```powershell
$body = @{
    idInstance = '1101000001'
    apiTokenInstance = 'your_token'
} | ConvertTo-Json

Invoke-RestMethod `
    -Method Post `
    -Uri 'http://localhost:8080/api/v1/state' `
    -ContentType 'application/json' `
    -Body $body
```

### Send Text Message

```powershell
$body = @{
    idInstance = '1101000001'
    apiTokenInstance = 'your_token'
    chatId = '79991234567@c.us'
    message = 'TEST ONLY'
    typingTime = 3000
} | ConvertTo-Json

Invoke-RestMethod `
    -Method Post `
    -Uri 'http://localhost:8080/api/v1/messages/text' `
    -ContentType 'application/json' `
    -Body $body
```

### Send File By URL

```powershell
$body = @{
    idInstance = '1101000001'
    apiTokenInstance = 'your_token'
    chatId = '79991234567@c.us'
    urlFile = 'https://example.com/files/report.pdf'
    fileName = 'report.pdf'
    caption = 'Test file'
} | ConvertTo-Json

Invoke-RestMethod `
    -Method Post `
    -Uri 'http://localhost:8080/api/v1/messages/file' `
    -ContentType 'application/json' `
    -Body $body
```

## WhatsApp Testing Notes

Чтобы не отправить сообщение случайно другому человеку, можно тестировать на своем личном номере.

Пример `chatId` для личного номера:

```text
79991234567@c.us
```

`quotedMessageId` — ID сообщения, на которое отправляется ответ. Для обычного теста можно не передавать.

`typingTime` — время в миллисекундах, в течение которого Green API эмулирует набор текста перед отправкой. В чате с самим собой индикатор `печатает...` может визуально не отображаться.

## Validation Rules

### `sendMessage`

- `chatId` обязателен
- `message` обязателен
- `message` до `20000` символов
- `typingTime` от `1000` до `20000`

### `sendFileByUrl`

- `chatId` обязателен
- `urlFile` обязателен и должен быть `https`
- `fileName` обязателен и должен содержать расширение
- `caption` до `20000` символов
- `typingType`, если передан, должен быть `recording`

## Tests

Локально:

```powershell
go test ./...
```

В текущей Windows-среде package tests проходят, но cleanup временного Go cache может завершаться с `Access is denied` уже после выполнения самих тестов.

## Notes

- Credentials не хранятся на сервере и передаются вручную в каждом запросе
- UI предназначен для ручного тестирования Green API
- Backend не зависит от дефолтных `idInstance` и `apiTokenInstance`
