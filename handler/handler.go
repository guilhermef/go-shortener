package handler

import (
	"gopkg.in/redis.v3"
	"log"
	"net/http"
	"time"
)

type RedirectHandler struct {
	Client *redis.Client
	Logger *log.Logger
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logEntry := time.Now().UTC().String()

	redirect, _ := h.Client.Get("go-shortener:" + req.RequestURI).Result()
	if redirect == "" {
		logEntry += " 404 " + req.RequestURI
		http.NotFound(w, req)
	} else {
		logEntry += " 301 " + req.RequestURI + " " + redirect
		h.Client.Incr("go-shortener-count:" + req.RequestURI)
		http.Redirect(w, req, redirect, 301)
	}

	h.Logger.Print(logEntry)
}
