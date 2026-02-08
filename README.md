# Health Check Notification System

Sistem untuk memantau kesehatan service dan mengirimkan notifikasi ketika service DOWN.

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- WAHA instance (for WhatsApp notifications)

### Run with Docker Compose

```bash
# Clone and enter directory
cd healt-check-service

# Start all services
docker-compose up -d

# View logs
docker-compose logs -f app
```

### Run Locally (Development)

```bash
# Start MongoDB & Redis
docker-compose up -d mongo redis

# Copy and configure .env
cp .env.example .env

# Run the application
go run cmd/main.go
```

## 🔔 Notification Verification

To ensure your notification channels are correctly configured, you can use the provided verification scripts.

### 1. Get Telegram Chat ID

If you don't know your Telegram Chat ID:

1. Start a chat with your bot.
2. Run the helper script:

```bash
go run cmd/get_telegram_chat_id/main.go
```

### 2. Verify Notifications

To send a test message to all configured channels:

1. Open `cmd/verify_notification/main.go` and update the recipient IDs (e.g., `REPLACE_WITH_CHAT_ID`).
2. Run the verification script:

```bash
go run cmd/verify_notification/main.go
```

## 📊 Dashboard

Access the dashboard at: `http://localhost:8080/dashboard`

Features:

- View service health status
- Add/Edit/Delete services to monitor
- Manage notification recipients
- View health check logs

## 📡 API Endpoints

### Services

| Method | Endpoint            | Description        |
| ------ | ------------------- | ------------------ |
| GET    | `/api/services`     | List all services  |
| GET    | `/api/services/:id` | Get service by ID  |
| POST   | `/api/services`     | Create new service |
| PUT    | `/api/services/:id` | Update service     |
| DELETE | `/api/services/:id` | Delete service     |

### Recipients

| Method | Endpoint              | Description         |
| ------ | --------------------- | ------------------- |
| GET    | `/api/recipients`     | List all recipients |
| POST   | `/api/recipients`     | Add new recipient   |
| PUT    | `/api/recipients/:id` | Update recipient    |
| DELETE | `/api/recipients/:id` | Delete recipient    |

### Health Logs

| Method | Endpoint                      | Description              |
| ------ | ----------------------------- | ------------------------ |
| GET    | `/api/health-logs`            | Get recent logs          |
| GET    | `/api/health-logs/:serviceId` | Get logs by service      |
| GET    | `/api/dashboard/stats`        | Get dashboard statistics |

## ⚙️ Configuration

Environment variables (`.env`):

```env
# MongoDB
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=health_check_db

# Redis (for AsyncQ)
REDIS_ADDR=localhost:6379

# WAHA
WAHA_API_URL=http://localhost:3000
WAHA_SESSION=default
WAHA_API_KEY=your_api_key

# App
APP_PORT=8080
DEFAULT_CHECK_INTERVAL=60
DEFAULT_RETRY_INTERVAL=180
HTTP_TIMEOUT=10

# Telegram
TELEGRAM_BOT_TOKEN=your_telegram_bot_token

# Discord
DISCORD_BOT_TOKEN=your_discord_bot_token
```

## 🏗️ Architecture

```
┌─────────────────┐     ┌─────────────────┐
│   Scheduler     │────▶│  Health Checker │
└─────────────────┘     └────────┬────────┘
                                 │
                    ┌────────────┴────────────┐
                    ▼                         ▼
            ┌─────────────┐           ┌─────────────┐
            │   MongoDB   │           │    AsyncQ   │
            │  (logs)     │           │   (Redis)   │
            └─────────────┘           └──────┬──────┘
                                             │
                       ┌────────────────────┼────────────────────┐
                       ▼                    ▼                    ▼
                ┌─────────────┐      ┌─────────────┐      ┌─────────────┐
                │    WAHA     │      │  Telegram   │      │   Discord   │
                │  (WhatsApp) │      │    Bot      │      │     Bot     │
                └─────────────┘      └─────────────┘      └─────────────┘
```

## 📁 Project Structure

```
healt-check-service/
├── cmd/main.go                 # Entry point
├── config/                     # Configuration
├── api/                        # REST API handlers
├── internal/
│   ├── domain/                 # Domain models
│   ├── repository/             # Data access
│   ├── usecase/                # Business logic
│   ├── queue/                  # AsyncQ
│   └── scheduler/              # Health check scheduler
├── infrastructure/waha/        # WAHA client
├── infrastructure/telegram/    # Telegram client
├── infrastructure/discord/     # Discord client
├── web/dashboard/              # Dashboard UI
├── docker-compose.yml
└── Dockerfile
```
