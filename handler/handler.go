package handler

import (
	"gopkg.in/redis.v3"
	"io"
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

	if req.RequestURI == "/healthcheck" {
		err := h.Client.Ping().Err()
		if err != nil {
			http.Error(w, "Could not connect to REDIS server", 502)
			return
		}
		logEntry += " 200 " + req.RequestURI
		h.Logger.Print(logEntry)
		io.WriteString(w, "WORKING\n")
		return
	}

	redirect, err := h.Client.Get("go-shortener:" + req.RequestURI).Result()
	if redirect == "" || err != nil {
		logEntry += " 404 " + req.RequestURI
		http.NotFound(w, req)
	} else {
		logEntry += " 301 " + req.RequestURI + " " + redirect
		h.Client.Incr("go-shortener-count:" + req.RequestURI)
		http.Redirect(w, req, redirect, 301)
	}

	h.Logger.Print(logEntry)
}
