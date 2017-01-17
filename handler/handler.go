package handler

import (
	"gopkg.in/redis.v3"
	"io"
	"log"
	"net/http"
	"time"
	"strconv"
	"os"
)

type RedirectHandler struct {
	Client *redis.Client
	Logger *log.Logger
	Extra  Extra
}

type Extra struct {
  RedirectHost string
  RedirectCode int
}

func getExtra(h *Extra) Extra{
	redirectCode := 302

	if os.Getenv("REDIRECT_CODE") != "" {
		redirectCode, _ = strconv.Atoi(os.Getenv("REDIRECT_CODE"))
	} else if(h.RedirectCode != 0) {
		redirectCode = h.RedirectCode
	}

	redirectHost := os.Getenv("REDIRECT_HOST")
	if redirectHost == "" {
		redirectHost = h.RedirectHost
	}

	return Extra{RedirectCode: redirectCode, RedirectHost: redirectHost}
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
		e := getExtra(&h.Extra)

		if e.RedirectHost != "" {
			logEntry += " " + strconv.Itoa(e.RedirectCode) + " " + req.RequestURI + " " + e.RedirectHost
			http.Redirect(w, req, e.RedirectHost, e.RedirectCode)
		} else {
			logEntry += " 404 " + req.RequestURI
			http.NotFound(w, req)
		}
	} else {
		logEntry += " 301 " + req.RequestURI + " " + redirect
		h.Client.Incr("go-shortener-count:" + req.RequestURI)
		http.Redirect(w, req, redirect, 301)
	}

	h.Logger.Print(logEntry)
}
