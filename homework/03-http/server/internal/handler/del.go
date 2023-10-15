package handler

import (
	"net/http"
	"os"
)

func (h Handler) del(w http.ResponseWriter, r *http.Request) {
	filePath := h.cfg.WorkingDirectory + r.URL.Path
	err := delFile(filePath, r)
	if err != nil {
		w.Header().Set("Server", h.cfg.Domain)
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	w.Header().Set("Server", h.cfg.Domain)
	w.WriteHeader(http.StatusOK)
	return

}

func delFile(filePath string, r *http.Request) error {
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if fileInfo.IsDir() {
			if r.Header.Get("Remove-Directory") != "True" {
				return ErrIsDir
			}
			err := os.RemoveAll(filePath)
			if err != nil {
				return err
			}
			return nil
		}
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}
