package main

import (
	"encoding/gob"
	"os"
	"sync"
	"time"

	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/data"
	"github.com/dubey22rohit/heyyy_yo_backend/auth/internal/mailer"
	"github.com/dubey22rohit/heyyy_yo_backend/pkg/jsonlog"
	"github.com/redis/go-redis/v9"
)

const version = "1.0.0"

type Config struct {
	Port  int
	Debug bool
	Env   string
	DB    struct {
		DSN          string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}
	RedisURL string
	SMTP     struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
	FrontendURL     string
	TokenExpiration struct {
		DurationString string
		Duration       time.Duration
	}
	Secret struct {
		HMC               string
		SecretKey         []byte
		SessionExpiration time.Duration
	}
}

type Application struct {
	Config      Config
	Logger      *jsonlog.Logger
	RedisClient *redis.Client
	wg          sync.WaitGroup
	Models      data.Models
	Mailer      mailer.Mailer
}

func main() {
	// The only data we will encrypt in the cookies is the UserID type. We need to register this here
	gob.Register(&data.UserID{})

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := updateConfigWithEnvVariables()
	if err != nil {
		logger.PrintFatal(err, nil, cfg.Debug)
	}

	db, err := openDB(*cfg)
	if err != nil {
		logger.PrintFatal(err, nil, cfg.Debug)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil, cfg.Debug)

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		logger.PrintFatal(err, nil, cfg.Debug)
	}

	client := redis.NewClient(opt)

	logger.PrintInfo("redis connection pool established", nil, cfg.Debug)

	app := &Application{
		Config:      *cfg,
		Logger:      logger,
		RedisClient: client,
		Models:      data.NewModels(db),
		Mailer:      mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil, cfg.Debug)
	}

}
