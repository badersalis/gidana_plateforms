package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	Port            string
	DatabaseURL     string
	DBPath          string
	JWTSecret       string
	JWTExpiryHours  int
	UploadDir       string
	MaxUploadSizeMB int64
	UseFirebase      bool
	FirebaseCredJSON string
	FirebaseCredPath string
	FirebaseBucket   string
	AllowedOrigins  string
}

var App *Config

func Load() {
	_ = godotenv.Load()

	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "72"))
	maxUpload, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE_MB", "5"), 10, 64)
	useFirebase, _ := strconv.ParseBool(getEnv("USE_FIREBASE", "false"))

	App = &Config{
		AppEnv:          getEnv("APP_ENV", "development"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		DBPath:          getEnv("DB_PATH", "gidana_dev.db"),
		JWTSecret:       getEnv("JWT_SECRET", "dev-secret-change-me"),
		JWTExpiryHours:  jwtExpiry,
		UploadDir:       getEnv("UPLOAD_DIR", "./uploads/properties"),
		MaxUploadSizeMB: maxUpload,
		UseFirebase:      useFirebase,
		FirebaseCredJSON: getEnv("FIREBASE_CREDENTIALS_JSON", ""),
		FirebaseCredPath: getEnv("FIREBASE_CREDENTIALS_PATH", ""),
		FirebaseBucket:   getEnv("FIREBASE_BUCKET", ""),
		AllowedOrigins:  getEnv("ALLOWED_ORIGINS", "*"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
