package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func (h Handler) put(w http.ResponseWriter, r *http.Request) {
	filePath := h.cfg.WorkingDirectory + r.URL.Path
	//log.Println(r.ContentLength)
	//log.Printf("Request: %+v\n", r)
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	err := createFilePut(filePath, r)
	if err != nil {
		w.Header().Set("Server", h.cfg.Domain)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.Header().Set("Server", h.cfg.Domain)
	w.WriteHeader(http.StatusCreated)
	return

}

func createFilePut(filePath string, r *http.Request) error {
	fileInfo, err := os.Stat(filePath)
	if err == nil {
		if fileInfo.IsDir() {

			return fmt.Errorf("is a directory")
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("Could not create file")
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		return fmt.Errorf("Could not save file")
	}

	return nil
}
