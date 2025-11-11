# 🔗 URL Shortener

Высокопроизводительный сервис сокращения URL, написанный на Go с использованием PostgreSQL и Redis.

## 🌟 Возможности

- ✅ Сокращение длинных URL в короткие коды
- ✅ Перенаправление по короткому коду на оригинальный URL
- ✅ Кэширование результатов в Redis для быстрого доступа
- ✅ Подсчёт количества переходов для каждого URL
- ✅ Валидация URL с автоматическим добавлением схемы (https://)
- ✅ REST API с поддержкой CORS
- ✅ Логирование с помощью zerolog
- ✅ Миграции базы данных

## 📋 Требования

- **Go** 1.19+
- **PostgreSQL** 12+
- **Redis** 6+
- **Docker** (опционально, для контейнеризации)

## 🚀 Быстрый старт

### 1️⃣ Клонирование репозитория

```bash
git clone https://github.com/massonsky/url-shortener
cd url-shortener
```

### 2️⃣ Установка зависимостей

```bash
go mod download
go mod tidy
```

### 3️⃣ Конфигурация окружения

Создайте файл `.env` в корне проекта:

```env
# Сервер
SERVER_PORT=8080
APP_ENV=development

# PostgreSQL
DATABASE_URL=postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable

# Redis
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
```

### 4️⃣ Запуск сервисов с Docker (опционально)

```bash
docker-compose up -d
```

Или установите PostgreSQL и Redis вручную.

### 5️⃣ Применение миграций

```bash
go run -tags 'postgres' ./cmd/migrate/main.go up
```

### 6️⃣ Запуск сервера

```bash
go build -o url-shortener ./cmd/server
./url-shortener
```

Сервер будет доступен по адресу: `http://localhost:8080`

## 📚 API Документация

### Сокращение URL

**Запрос:**
```http
POST /shorten
Content-Type: application/json

{
  "url": "https://www.example.com/very/long/url/path"
}
```

**Ответ (200 OK):**
```json
{
  "short_url": "http://localhost:8080/abc123"
}
```

**Примеры:**

```powershell
# PowerShell
$body = @{ url = "https://www.golang.org" } | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/shorten" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body
```

```bash
# Bash/curl
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.golang.org"}'
```

### Перенаправление по короткому коду

**Запрос:**
```http
GET /{code}
```

**Ответ:**
- **302 Found** — перенаправление на оригинальный URL
- **404 Not Found** — короткий код не найден

**Пример:**
```bash
curl -L http://localhost:8080/1
# Перенаправит на https://www.golang.org
```

## 🏗️ Архитектура

```
url-shortener/
├── cmd/
│   └── server/
│       └── main.go           # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go         # Конфигурация из .env
│   ├── domain/
│   │   └── url.go            # Доменная модель
│   ├── repository/
│   │   ├── postgres/
│   │   │   └── url_repo.go   # Работа с БД
│   │   └── redis/
│   │       └── client.go     # Redis клиент
│   ├── service/
│   │   └── url_service.go    # Бизнес-логика
│   ├── transport/
│   │   └── _http/
│   │       ├── handler.go    # HTTP обработчики
│   │       ├── middleware.go # HTTP middleware
│   │       └── server.go     # HTTP сервер
│   └── utils/
│       └── base62.go         # Кодирование в base62
├── migration/
│   └── 000001_create_urls_table.up.sql  # Миграции БД
├── docker-compose.yml        # Docker конфигурация
└── README.md                 # Этот файл
```

## 🔄 Как это работает

### Процесс сокращения URL

1. Пользователь отправляет длинный URL в `/shorten`
2. Handler валидирует и нормализует URL (добавляет https:// если нужно)
3. Service парсит URL и проверяет корректность
4. Repository сохраняет URL в PostgreSQL и получает `id`
5. `id` кодируется в base62 → получается короткий код
6. Короткий код сохраняется в БД и кэшируется в Redis
7. Возвращается ссылка `http://localhost:8080/{shortCode}`

### Процесс перенаправления

1. Пользователь переходит по ссылке `http://localhost:8080/{shortCode}`
2. Handler получает `shortCode` из URL
3. Service сначала ищет URL в Redis кэше (быстро ⚡)
4. Если в кэше нет → ищет в PostgreSQL
5. Фоновый процесс обновляет счётчик переходов
6. Возвращается 302 redirect на оригинальный URL

## 🗄️ Схема БД

```sql
CREATE TABLE urls (
  id SERIAL PRIMARY KEY,
  original_url TEXT NOT NULL,
  short_code VARCHAR(255) UNIQUE NOT NULL,
  click_count INT DEFAULT 0,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_original_url ON urls(original_url);
```

## 🛠️ Используемые технологии

| Компонент | Пакет | Назначение |
|-----------|-------|-----------|
| **Веб-фреймворк** | [chi](https://github.com/go-chi/chi) | Маршрутизация HTTP |
| **CORS** | [chi/cors](https://github.com/go-chi/cors) | Поддержка CORS |
| **Валидация** | [validator](https://github.com/go-playground/validator) | Валидация данных |
| **БД** | [pgx](https://github.com/jackc/pgx) | PostgreSQL драйвер |
| **Кэш** | [redis](https://github.com/redis/go-redis) | Redis клиент |
| **Миграции** | [migrate](https://github.com/golang-migrate/migrate) | Управление миграциями |
| **Логирование** | [zerolog](https://github.com/rs/zerolog) | Структурированное логирование |

## 📝 Примеры использования

### Сокращение URL с пробелами

```powershell
$body = @{ url = " https://www.golang.org   " } | ConvertTo-Json
Invoke-RestMethod -Uri "http://localhost:8080/shorten" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body
```

**Результат:**
```json
{
  "short_url": "http://localhost:8080/1"
}
```

### Перенаправление

```powershell
$response = Invoke-WebRequest -Uri "http://localhost:8080/1" -MaximumRedirection 0 -ErrorAction SilentlyContinue
$response.Headers.Location  # https://www.golang.org
```

### Сокращение URL без схемы

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "example.com"}'
# Автоматически станет https://example.com
```

## 🐛 Отладка

### Логирование

Приложение использует **zerolog** для структурированного логирования:

```go
// В коде
log.Error().Err(err).Msg("Failed to shorten URL")

// В консоли
{"level":"error","error":"invalid URL","time":"2024-11-07T10:30:45Z","message":"Failed to shorten URL"}
```

### Проверка подключения к БД

```bash
# PostgreSQL
psql postgresql://postgres:postgres@localhost:5432/shortener

# Redis
redis-cli ping
# PONG
```

## 📊 Производительность

- **Создание короткого URL**: ~10ms (с БД)
- **Перенаправление (из кэша)**: ~1ms
- **Перенаправление (из БД)**: ~5-10ms

## 📄 Лицензия

MIT License — см. файл [LICENSE](LICENSE)
