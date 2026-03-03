package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort                  string
	DBHost                    string
	DBPort                    string
	DBUser                    string
	DBPassword                string
	DBName                    string
	JWTSecret                 string
	JWTExpireMinutes          int
	RefreshTokenExpireMinutes int
}

func Load() (*Config, error) {
	httpPort := getEnvOrDefault("HTTP_PORT", "8080")
	dbHost := getEnvOrDefault("POSTGRES_HOST", "localhost")
	dbPort := getEnvOrDefault("POSTGRES_PORT", "5432")
	dbUser := getEnvOrDefault("POSTGRES_USER", "auth_user")
	dbPassword := getEnvOrDefault("POSTGRES_PASSWORD", "auth_password")
	dbName := getEnvOrDefault("POSTGRES_DB", "auth_db")
	jwtSecret := getEnvOrDefault("JWT_SECRET", "secret")
	jwtExpireMinutes := parseEnvInt("JWT_EXPIRE_MINUTES", 30)
	refreshTokenExpireMinutes := parseEnvInt("REFRESH_TOKEN_EXPIRE_MINUTES", 60)

	return &Config{
		HTTPPort:                  httpPort,
		DBHost:                    dbHost,
		DBPort:                    dbPort,
		DBUser:                    dbUser,
		DBPassword:                dbPassword,
		DBName:                    dbName,
		JWTSecret:                 jwtSecret,
		JWTExpireMinutes:          jwtExpireMinutes,
		RefreshTokenExpireMinutes: refreshTokenExpireMinutes,
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func parseEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
