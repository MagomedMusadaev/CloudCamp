# Load Balancer с Rate Limiting

Высокопроизводительный HTTP балансировщик нагрузки с поддержкой rate limiting, написанный на Go. Проект включает различные стратегии балансировки, мониторинг здоровья бэкендов и управление ограничениями для клиентов.

## Основные возможности

- **Балансировка нагрузки**:
  - Round Robin (циклическое распределение)
  - Least Connections (наименьшее количество соединений)
  - Random (случайное распределение)
- **Rate Limiting**:
  - Token Bucket алгоритм
  - Поддержка глобальных и клиентских лимитов
  - Настраиваемые периоды и лимиты
- **Мониторинг**:
  - Health checks для бэкендов
  - Автоматическое исключение недоступных серверов
- **Управление**:
  - CRUD API для управления клиентами
  - Конфигурация через YAML
  - Graceful shutdown

## Требования

- Go 1.21+
- Docker и Docker Compose (для запуска в контейнерах)

## Установка и запуск

### Локальный запуск

```bash
# Клонирование репозитория
git clone https://github.com/MagomedMusadaev/CloudCamp
cd CloudCamp

# Установка зависимостей
go mod download

# Запуск приложения
go run cmd/balancer/main.go
```

### Запуск через Docker Compose

```bash
# Сборка и запуск контейнеров
docker-compose up --build
```

## Конфигурация

Настройки приложения находятся в `configs/config.yaml`:

```yaml
env: development  # Доступные окружения: development, production, test

server:
  port: 8080

balancer:
  backends:
    - "http://backend1:8081"
    - "http://backend2:8082"
    - "http://backend3:8083"
  strategy: round-robin         # Доступные стратегии: round-robin, random, least-connections

health_checker:
  enabled: true
  interval: 15s                 # Интервал проверки
  path: "/health"               # Путь для проверки здоровья

rate_limiter:
  enabled: true
  interval: 10s                 # Интервал обновления токенов (global)
  rate: 35                      # Глобальный лимит (всего токенов)
  period: 2m                    # Период для глобального лимита

  clients:
    client1:
      rate: 7
      period: 1m

    client2:
      rate: 23
      period: 2m

log:
  file_path: "./logs/app.log"   # Путь к файлу логов
  dir: "./logs"                 # Директория для логов
```

## API Endpoints

### Прокси-сервер

Основной функционал балансировщика доступен через прокси-эндпоинт, который обрабатывает все входящие HTTP-запросы и перенаправляет их на доступные бэкенды.

Request:
- Поддерживаются все HTTP методы (GET, POST, PUT, DELETE, etc.)
- Все заголовки и тело запроса передаются на бэкенд без изменений
- Для авторизованных клиентов добавить заголовок: X-Client-ID: <client_id>

Response:
- Статус и заголовки от бэкенда передаются клиенту
- Добавляются служебные заголовки для отладки

Ошибки:
- 429 Too Many Requests: превышен лимит запросов
- 502 Bad Gateway: ошибка взаимодействия с бэкендом
- 503 Service Unavailable: нет доступных бэкендов

Примеры запросов:

```http
# Базовый запрос
curl http://localhost:8080

# Запрос с указанием клиента
curl -H "X-Client-ID: client1" http://localhost:8080
```

### Управление клиентами

#### Создание клиента
```http request
POST /clients
Content-Type: application/json

Request:
{
    "client_id": "user1",
    "rate": 10,
    "period": "1m"
}

Response 201:
{
    "message": "Client created successfully"
}

Response 405:
{
    "code": "405",
    "message": "Method not allowed"
 
}

Response 400:
{
    "code": "400",
    "message": "Invalid request body", "Invalid client parameters", "Invalid period format. Example: '1s', '500ms', '2m'"
}

Response 500:
{
    "code": "500",
    "message": "Failed to set client limit"
}
```

#### Удаление клиента
```http
DELETE /clients?clientId=user1

Response 200:
{
    "code": "200",
    "message": "Client deleted successfully"
}

Response 405:
{
    "code": "405",
    "message": "Method not allowed"
}

Response 400:
{
    "code": "400",
    "message": "Client ID is required"
}
```

## Тестирование

### Запуск тестов

```bash
# Запуск тестов
go test ./tests -bench=. -race -v 
```

### Нагрузочное тестирование

```bash
# Пример использования Apache Bench
ab -n 5000 -c 1000 http://localhost:8080/
```

## Архитектура

Проект следует принципам чистой архитектуры и разделен на следующие основные компоненты:

- `cmd/` - точка входа приложения
- `internal/`
  - `app/` - конфигурация сервера и маршрутизация
  - `background/` - фоновые джобы refill и health
  - `balancer/` - реализации стратегий балансировки
  - `config/` - конфигурация приложения и сборка параметров
  - `domain/` - бизнес-логика и интерфейсы
  - `handler/` - HTTP обработчики
  - `limiter/` - реализация rate limiting
  - `logger/` - настройка логирования
- `pkg/` - общие утилиты
- `tests/` - тесты

## Комментирование кода

Весь код проекта тщательно прокомментирован для удобства изучения и чтения. Каждый пакет, структура, интерфейс и функция содержат подробные комментарии, объясняющие их назначение и принцип работы.

## Логирование

Логи сохраняются в директории `logs/`. Уровень логирования можно настроить в конфигурационном файле.