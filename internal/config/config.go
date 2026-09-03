package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Env                string
	DBDriver           string // "postgres" or "memory"
	DatabaseURL        string
	JWTSecret          string
	JWTExpirationHours int
	AdminEmail         string
	AdminPassword      string
	AllowedOrigins     []string

	// SMTP / Email Settings
	SMTPHost          string
	SMTPPort          int
	SMTPUser          string
	SMTPPassword      string
	SMTPFrom          string
	NotificationEmail string // Destination email for appointment notifications
}

func Load() *Config {
	// Attempt to load .env file if it exists
	if err := godotenv.Load(); err != nil {
		slog.Debug(".env file not found or could not be loaded, using environment variables")
	}

	jwtExpHours, _ := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "24"))
	if jwtExpHours <= 0 {
		jwtExpHours = 24
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if smtpPort <= 0 {
		smtpPort = 587
	}

	originsRaw := getEnv("CORS_ALLOWED_ORIGINS", "*")
	origins := strings.Split(originsRaw, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}

	adminEmail := getEnv("ADMIN_EMAIL", "admin@mariadiezma.com")

	return &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		DBDriver:           getEnv("DB_DRIVER", "postgres"),
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/mariadiezma?sslmode=disable"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-jwt-key-change-me-in-production"),
		JWTExpirationHours: jwtExpHours,
		AdminEmail:         adminEmail,
		AdminPassword:      getEnv("ADMIN_PASSWORD", "AdminPass123!"),
		AllowedOrigins:     origins,

		SMTPHost:          getEnv("SMTP_HOST", ""),
		SMTPPort:          smtpPort,
		SMTPUser:          getEnv("SMTP_USER", ""),
		SMTPPassword:      getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:          getEnv("SMTP_FROM", "no-reply@mariadiezma.com"),
		NotificationEmail: getEnv("NOTIFICATION_EMAIL", adminEmail),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}
