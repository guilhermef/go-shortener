package handler

import (
	"net/http"
  "gopkg.in/redis.v3"
)

type RedirectHandler struct {
	client *redis.Client
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  redirect, _ := h.client.Get("go-shortner:" + req.RequestURI).Result()
  http.Redirect(w, req, redirect, 301)
}
