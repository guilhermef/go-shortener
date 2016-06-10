package config

import (
	"errors"
	"gopkg.in/redis.v3"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type readSettings struct {
	LogPath   string
	Port      string
	RedisHost string
}

type Config struct {
	Logger      *log.Logger
	Port        string
	RedisClient *redis.Client
}

func try(value string, def string) string {
	if value != "" {
		return value
	}

	return def
}

func getRedisOptions(settings *readSettings) *redis.Options {
	return &redis.Options{
		Addr:     try(settings.RedisHost, "localhost:6379"),
		Password: "", // no password set
		DB:       0,  // use default DB
	}
}

func getLoggerOutput(settings *readSettings) io.Writer {
	if settings.LogPath == "" {
		return os.Stdout
	}

	file, _ := os.Create(settings.LogPath)
	return file
}

func fetchSettings() (*readSettings, error) {
	s := readSettings{LogPath: "", Port: "1234", RedisHost: "localhost:6379"}
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
		Port:        try(s.Port, "1234"),
		RedisClient: redisClient}, err
}
