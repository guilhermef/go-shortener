package handler

import (
	"net/http"
  "log"
	"gopkg.in/redis.v3"
  "time"
)

type RedirectHandler struct {
	client *redis.Client
  logger *log.Logger
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  logEntry := time.Now().UTC().String()

	redirect, _ := h.client.Get("go-shortener:" + req.RequestURI).Result()
	if redirect == "" {
    logEntry += " 404 " + req.RequestURI
		http.NotFound(w, req)
	} else {
    logEntry += " 301 " + req.RequestURI + " " + redirect
  	h.client.Incr("go-shortener-count:" + req.RequestURI)
  	http.Redirect(w, req, redirect, 301)
  }

  h.logger.Print(logEntry)
}
