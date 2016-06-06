package main

import (
	"github.com/guilhermef/go-shortener/handler"
	//"gopkg.in/redis.v3"
	"log"
	"net/http"
	"os"
)

type config struct {
	logPath   string
	redisHost string
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
	config.logPath = getEnvOrDefault("LOG_PATH", "")
}

func main() {
	cfg := &config{}
	setConfig(cfg)
	initialize(cfg)
}

func initialize(cfg *config) {
	err := http.ListenAndServe(addr, &handler.RedirectHandler{})
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
