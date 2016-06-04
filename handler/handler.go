package handler

import {
    "gopkg.in/redis.v3"
    "net/http"
}

type RedirectHandler struct {
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "http://location.test", 301)
}
