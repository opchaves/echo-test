package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
	eLog "github.com/labstack/gommon/log"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

var doOnce sync.Once

type ctxKey int

const (
	CtxClaims ctxKey = iota
	CtxRefreshToken
	CtxVersion

	// EnvTest represents the test environment
	EnvTest string = "test"
)

var (
	Name    = getEnv("APP_NAME", "kommonei")
	Env     = getEnv("APP_ENV", "development")
	Host    = getEnv("HOST", "0.0.0.0")
	Port    = getEnv("PORT", "8080")
	Origins = getEnv("ORIGINS", "")

	IsDevelopment = Env == "development"
	IsProduction  = Env == "production"
	IsTest        = Env == "test"
	IsLocal       = IsDevelopment || IsTest

	DbHost      = getEnv("DB_HOST", "localhost")
	DbPort      = getEnv("DB_PORT", "5432")
	DbUser      = getEnv("DB_USER", "postgres")
	DbPassword  = getEnv("DB_PASSWORD", "secret")
	DbName      = getEnv("DB_NAME", "app_dev")
	DbSSLMode   = getEnv("DB_SSL_MODE", "disable")
	DatabaseURL = getDatabaseURL()

	JwtSecret        = getEnv("JWT_SECRET", "superSecret")
	JwtExpiry        = toDuration("JWT_EXPIRY", "1h")
	JwtRefreshExpiry = toDuration("JWT_REFRESH_EXPIRY", "72h")

	LoginTokenURL    = getEnv("LOGIN_TOKEN_URL", "http://localhost:8080/auth/token")
	LoginTokenLength = toInt("LOGIN_TOKEN_LENGTH", "16")
	LoginTokenExpiry = toDuration("LOGIN_TOKEN_EXPIRY", "10m")

	EmailSmtpHost = ""
	EmailSmtpPort = 0
	EmailSmtpUser = ""
	EmailSmptPass = ""
	EmailFromAddr = ""
	EmailFromName = ""

	LetterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	Root = filepath.Join(basepath, "..")

	LogLevel = getLogLevel()
)

func toDuration(envVar string, defaultVal string) time.Duration {
	val, err := time.ParseDuration(getEnv(envVar, defaultVal))
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", envVar, err)
	}
	return val
}

func toInt(envVar string, defaultVal string) int {
	val, err := strconv.Atoi(getEnv(envVar, defaultVal))
	if err != nil {
		log.Fatalf("Invalid value for %s: %s", envVar, err)
	}
	return val
}

func getEnv(name, defaultValue string) string {
	doOnce.Do(func() {
		slog.Default().Info("initializing config")
		InitConfig()
	})

	if value := os.Getenv(name); value != "" {
		return value
	}

	return defaultValue
}

func InitConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	env = strings.ToLower(env)

	log.Printf("Loading environment: %s", env)

	loadEnvVars(env)
}

func loadEnvVars(env string) {
	if os.Getenv("APP_ENV") == "production" {
		return
	}

	if !isValidEnv(env) {
		log.Fatalf("Invalid environment: %s", env)
	}

	filename := fmt.Sprintf("%s/../.env", basepath)
	slog.Default().Info("loading env variables.", "file", filename)

	err := godotenv.Load(filename)
	if err != nil {
		log.Printf("File %s not found. Using default values", filename)
	}
}

func isValidEnv(env string) bool {
	return env == "development" || env == "test" || env == "production"
}

func getDatabaseURL() string {
	if IsTest {
		return os.Getenv("DATABASE_URL_TEST")
	}

	return os.Getenv("DATABASE_URL")
}

// EnableTestEnv sets the APP_ENV to test
// This is useful for running tests
func EnableTestEnv() {
	if err := os.Setenv("APP_ENV", EnvTest); err != nil {
		panic(err)
	}
}

func getLogLevel() eLog.Lvl {
	if IsProduction {
		return eLog.INFO
	}

	return eLog.DEBUG
}
