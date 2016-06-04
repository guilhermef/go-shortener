package handler

import (
	"net/http"
  "gopkg.in/redis.v3"
)

type RedirectHandler struct {
	client *redis.Client
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  redirect, _ := h.client.Get("go-shortener:" + req.RequestURI).Result()
  h.client.Incr("go-shortener-count:" + req.RequestURI)
  http.Redirect(w, req, redirect, 301)
}
