package config

import (
	"errors"
	"gopkg.in/redis.v3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

type readSettings struct {
	LogPath   string
	Port      string
	RedisHost string
	RedisPass string
	RedisDB   int64
}

type Config struct {
	Logger      *log.Logger
	Port        string
	RedisClient *redis.Client
}

func try(env string, value string, def string) string {
	if env != "" {
		return env
	}
	if value != "" {
		return value
	}

	return def
}

func tryInt(env string, def int64) int64 {
	if env != "" {
		envInt, err := strconv.ParseInt(env, 10, 64)
		if err != nil {
			panic("Could not parse passed REDIS_DB environment variable")
		}
		return envInt
	}

	return def
}

func getRedisOptions(settings *readSettings) *redis.Options {
	return &redis.Options{
		Addr:     try(os.Getenv("REDIS_HOST"), settings.RedisHost, "127.0.0.1:6379"),
		Password: try(os.Getenv("REDIS_PASS"), settings.RedisPass, ""),
		DB:       tryInt(os.Getenv("REDIS_DB"), settings.RedisDB)}
}

func getLoggerOutput(settings *readSettings) io.Writer {
	logPath := try(os.Getenv("LOG_PATH"), settings.LogPath, "")
	if logPath == "" {
		return os.Stdout
	}

	file, _ := os.Create(settings.LogPath)
	return file
}

func fetchSettings() (*readSettings, error) {
	s := readSettings{LogPath: "", Port: "", RedisHost: ""}
	data, readErr := ioutil.ReadFile("./settings.yml")
	if readErr != nil {
		return &s, nil
	}

	yamlErr := yaml.Unmarshal([]byte(data), &s)
	if yamlErr != nil {
		return &s, errors.New("Unable to open settings.yml")
	}
	return &s, nil
}

func NewConfig() (*Config, error) {
	s, err := fetchSettings()

	if err != nil {
		return &Config{}, err
	}
	redisClient := redis.NewClient(getRedisOptions(s))
	logger := log.New(getLoggerOutput(s), "", 0)

	return &Config{
		Logger:      logger,
		Port:        try(os.Getenv("PORT"), s.Port, "1234"),
		RedisClient: redisClient}, err
}
