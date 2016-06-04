package handler

import "net/http"

type RedirectHandler struct {
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
}
