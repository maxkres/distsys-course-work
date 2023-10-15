package handler

import (
	"fmt"
	"net/http"
	"server/internal/config"
)

var (
	ErrFileNotFound = fmt.Errorf("file not found")
	ErrFileFound    = fmt.Errorf("file found")
	ErrIsDir        = fmt.Errorf("is a dir")
)

const maxRequestSize = 1024 * 1024 * 1024 // 1 GB

type Handler struct {
	cfg config.Config
}

func New(cfg config.Config) Handler {
	return Handler{
		cfg: cfg,
	}
}

func (h Handler) Controller(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		h.get(w, r)
	} else if r.Method == http.MethodPost {
		h.post(w, r)
	} else if r.Method == http.MethodPut {
		h.put(w, r)
	} else if r.Method == http.MethodDelete {
		h.del(w, r)
	}

}
