# Health Check Notification System

Health Check Service is a robust monitoring solution designed to track the availability of your HTTP services. It periodically checks the health of registered endpoints and sends real-time notifications when a service goes DOWN or recovers (UP).

**Key Features:**

- **Multi-Channel Notifications**: Support for WhatsApp (via WAHA), Telegram, and Discord.
- **Real-time Dashboard**: Web interface to view service status, logs, and manage configuration.
- **Flexible Scheduling**: Customizable check intervals and retry logic.
- **Persistent Logging**: Detailed history of health checks stored in MongoDB.

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- WAHA instance (for WhatsApp notifications)

### 📥 Installation

1.  **Clone the repository**:

    ```bash
    git clone https://github.com/RhmnDev14/health-check-service.git
    cd health-check-service
    ```

2.  **Install Dependencies**:

    ```bash
    go mod download
    ```

3.  **Setup Configuration**:
    ```bash
    cp .env.example .env
    # Edit .env with your database credentials and API tokens
    ```

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

```bash
healt-check-service/
├── cmd/                        # Command line entry points
│   ├── main.go                 # Main application server
│   ├── get_telegram_chat_id/   # Utility to retrieve Telegram Chat ID
│   └── verify_notification/    # Utility to test notification channels
├── config/                     # Configuration loader (env vars)
├── infrastructure/             # External service implementations
│   ├── discord/                # Discord client
│   ├── telegram/               # Telegram client
│   └── waha/                   # WhatsApp (WAHA) client
├── internal/                   # Private application code
│   ├── api/                    # HTTP Router and Server setup
│   ├── domain/                 # Domain entities and interfaces
│   ├── handler/                # HTTP Controllers/Handlers
│   ├── queue/                  # Async task queue (Redis)
│   ├── repository/             # Database persistence (MongoDB)
│   ├── scheduler/              # Cron-like scheduler for health checks
│   └── usecase/                # Business logic and application services
├── web/
│   └── dashboard/              # Frontend Dashboard (HTML/JS/CSS)
├── .env                        # Environment variables file
├── .env.example                # Example environment variables
├── docker-compose.yml          # Container orchestration config
├── Dockerfile                  # Application container definition
├── go.mod                      # Go module definitions
└── README.md                   # Project documentation
```
