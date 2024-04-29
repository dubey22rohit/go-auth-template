package main

import (
	"encoding/hex"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func updateConfigWithEnvVariables() (*Config, error) {
	err := godotenv.Load(".env", ".env.development")
	if err != nil {
		log.Fatal("error loading the .env file")
	}

	maxOpenConnsStr := os.Getenv("DB_MAX_OPEN_CONNS")
	maxOpenConns, err := strconv.Atoi(maxOpenConnsStr)
	if err != nil {
		log.Fatal(err)
	}

	maxIdleConnsStr := os.Getenv("DB_MAX_IDLE_CONNS")
	maxIdleConns, err := strconv.Atoi(maxIdleConnsStr)
	if err != nil {
		log.Fatal(err)
	}

	var cfg Config

	//Basic config
	flag.IntVar(&cfg.Port, "port", 8080, "Auth API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment(development | staging | production)")

	//Database config
	flag.StringVar(&cfg.DB.DSN, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", maxOpenConns, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MaxIdleConns, "db-max-idle-conns", maxIdleConns, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.DB.MaxIdleTime,
		"db-max-idle-time",
		os.Getenv("DB_MAX_IDLE_TIME"),
		"PostgreSQL max connection idle time",
	)

	//Redis config
	flag.StringVar(&cfg.RedisURL, "redis-url", os.Getenv("REDIS_URL"), "Redis URL")

	// Email
	emailPortStr := os.Getenv("EMAIL_SERVER_PORT")
	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&cfg.SMTP.Host, "smtp-host", os.Getenv("EMAIL_HOST_SERVER"), "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", emailPort, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", os.Getenv("EMAIL_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", os.Getenv("EMAIL_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", "hanzohasashi2212@gmail.com", "SMTP sender")

	// Frontend
	flag.StringVar(&cfg.FrontendURL, "frontend-url", os.Getenv("FRONTEND_URL"), "Frontend URL")

	// Token Expiration
	tokexpirationStr := os.Getenv("TOKEN_EXPIRATION")
	duration, err := time.ParseDuration(tokexpirationStr)
	if err != nil {
		return nil, err
	}
	cfg.TokenExpiration.DurationString = tokexpirationStr
	cfg.TokenExpiration.Duration = duration

	// Secret
	flag.StringVar(&cfg.Secret.HMC, "secret-key", os.Getenv("HMC_SECRET_KEY"), "HMC Secret Key")
	secretKey, err := hex.DecodeString(cfg.Secret.HMC)
	if err != nil {
		return nil, err
	}
	cfg.Secret.SecretKey = secretKey
	sessionDuration, err := time.ParseDuration(os.Getenv("SESSION_EXPIRATION"))
	if err != nil {
		return nil, err
	}
	cfg.Secret.SessionExpiration = sessionDuration

	flag.Parse()

	return &cfg, nil
}
