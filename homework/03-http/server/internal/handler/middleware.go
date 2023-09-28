package handler

import (
	"net/http"
)

func (h Handler) Host(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		host := req.Host
		if h.cfg.Domain != "" && h.cfg.Domain != host {
			w.Header().Set("Server", h.cfg.Domain)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, req)
	})
}
