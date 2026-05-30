package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv   string
	LogLevel string
	Port     string

	DatabaseURL          string
	DBMaxOpenConns       int
	DBMaxIdleConns       int
	DBConnMaxLifetimeMin int
	DBAutoMigrate        bool

	JWTSecret       string
	JWTExpiresHours int

	SeedDefaultAdmin     bool
	DefaultAdminEmail    string
	DefaultAdminPassword string

	RedisURL      string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	OTPTTLSeconds         int
	OTPMaxVerifyAttempts  int
	OTPRateLimitMax       int
	OTPRateLimitWindowMin int
	OTPDebugReturn        bool

	ResendAPIKey  string
	ResendBaseURL string
	EmailFrom     string

	StoragePath string
	MaxUploadMB int

	CORSAllowOrigins string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		AppEnv:   getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Port:     getEnv("PORT", "8080"),

		DatabaseURL:          getEnv("DATABASE_URL", ""),
		DBMaxOpenConns:       getEnvInt("DB_MAX_OPEN_CONNS", 20),
		DBMaxIdleConns:       getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetimeMin: getEnvInt("DB_CONN_MAX_LIFETIME_MIN", 30),
		DBAutoMigrate:        getEnvBool("DB_AUTO_MIGRATE", false),

		JWTSecret:       getEnv("JWT_SECRET", "change-me"),
		JWTExpiresHours: getEnvInt("JWT_EXPIRES_HOURS", 24),

		SeedDefaultAdmin:     getEnvBool("SEED_DEFAULT_ADMIN", true),
		DefaultAdminEmail:    getEnv("DEFAULT_ADMIN_EMAIL", "admin@pdfreader.local"),
		DefaultAdminPassword: getEnv("DEFAULT_ADMIN_PASSWORD", "PdfAdmin123@"),

		RedisURL:      getEnv("REDIS_URL", ""),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),

		OTPTTLSeconds:         getEnvInt("OTP_TTL_SECONDS", 600),
		OTPMaxVerifyAttempts:  getEnvInt("OTP_MAX_VERIFY_ATTEMPTS", 5),
		OTPRateLimitMax:       getEnvInt("OTP_RATE_LIMIT_MAX", 3),
		OTPRateLimitWindowMin: getEnvInt("OTP_RATE_LIMIT_WINDOW_MIN", 15),
		OTPDebugReturn:        getEnvBool("OTP_DEBUG_RETURN", false),

		ResendAPIKey:  getEnv("RESEND_API_KEY", ""),
		ResendBaseURL: getEnv("RESEND_BASE_URL", "https://api.resend.com"),
		EmailFrom:     getEnv("EMAIL_FROM", ""),

		StoragePath: getEnv("STORAGE_PATH", "./storage"),
		MaxUploadMB: getEnvInt("MAX_UPLOAD_MB", 25),

		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
