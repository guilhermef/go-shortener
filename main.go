package main

import (
	"github.com/guilhermef/go-shortener/handler"
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"os"
)

type config struct {
	logPath   string
	redisHost string
	port      string
}

const addr = "localhost:12345"

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func setConfig(config *config) {
	config.redisHost = getEnvOrDefault("REDIS_HOST", "localhost:6379")
	config.port = getEnvOrDefault("PORT", "1234")
	config.logPath = getEnvOrDefault("LOG_PATH", "")
}

func getLogger(cfg *config) *log.Logger {
	if cfg.logPath == "" {
		return log.New(os.Stdout, "", 0)
	}

	file, _ := os.Create(cfg.logPath)
	return log.New(file, "", 0)
}

func getRedisClient(cfg *config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	cfg := &config{}
	setConfig(cfg)
	initialize(cfg)
}

func initialize(cfg *config) {
	logger := getLogger(cfg)
	redisClient := getRedisClient(cfg)

	err := http.ListenAndServe("localhost:"+cfg.port, &handler.RedirectHandler{Client: redisClient, Logger: logger})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
