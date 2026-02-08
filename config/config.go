package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	// MongoDB
	MongoURI      string
	MongoDatabase string

	// Redis
	RedisAddr     string
	RedisPassword string

	// WAHA
	WahaAPIURL  string
	WahaSession string
	WahaAPIKey  string

	// Telegram
	TelegramBotToken string

	// Discord
	DiscordBotToken string

	// Health Check
	DefaultCheckInterval int
	DefaultRetryInterval int
	HTTPTimeout          int

	// App
	AppEnv   string
	LogLevel string
	AppPort  string
}

var AppConfig *Config

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Warn("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		// MongoDB
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "health_check_db"),

		// Redis
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),

		// WAHA
		WahaAPIURL:  getEnv("WAHA_API_URL", "http://localhost:3000"),
		WahaSession: getEnv("WAHA_SESSION", "default"),
		WahaAPIKey:  getEnv("WAHA_API_KEY", ""),

		// Telegram
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),

		// Discord
		DiscordBotToken: getEnv("DISCORD_BOT_TOKEN", ""),

		// Health Check
		DefaultCheckInterval: getEnvAsInt("DEFAULT_CHECK_INTERVAL", 60),
		DefaultRetryInterval: getEnvAsInt("DEFAULT_RETRY_INTERVAL", 180),
		HTTPTimeout:          getEnvAsInt("HTTP_TIMEOUT", 10),

		// App
		AppEnv:   getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "debug"),
		AppPort:  getEnv("APP_PORT", "8080"),
	}

	// Setup logrus
	level, err := logrus.ParseLevel(AppConfig.LogLevel)
	if err != nil {
		level = logrus.DebugLevel
	}
	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return AppConfig
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
