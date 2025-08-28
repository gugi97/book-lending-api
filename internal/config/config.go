package config

import (
    "os"
)

// Config holds all of the configuration values for the application.  Each
// nested struct corresponds to a portion of the configuration.  Values
// are pulled from environment variables with sensible defaults when
// unavailable.  No local .env loading is performed to avoid mixing
// development and production settings implicitly.
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    JWT      JWTConfig
}

// ServerConfig controls the HTTP server.
type ServerConfig struct {
    Port string
}

// DatabaseConfig describes the connection to MySQL.
type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
}

// JWTConfig holds the signing key used for JSON web tokens.
type JWTConfig struct {
    Secret string
}

// Load reads configuration from the environment.  If no value is set for
// a given variable the provided default will be used instead.  This
// function should be called as early as possible in program startup.
func Load() *Config {
    return &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "db"),
            Port:     getEnv("DB_PORT", "3306"),
            User:     getEnv("DB_USER", "root"),
            Password: getEnv("DB_PASSWORD", "password"),
            Name:     getEnv("DB_NAME", "book_lending"),
        },
        JWT: JWTConfig{
            Secret: getEnv("JWT_SECRET", "supersecretkey"),
        },
    }
}

// getEnv returns the value of the environment variable if set or the
// provided default otherwise.
func getEnv(key, defaultValue string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultValue
}