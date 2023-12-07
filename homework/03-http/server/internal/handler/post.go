package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func (h Handler) post(w http.ResponseWriter, r *http.Request) {
	filePath := h.cfg.WorkingDirectory + r.URL.Path
	//log.Printf("Request: %+v\n", r)
	r.Body = http.MaxBytesReader(w, r.Body, maxRequestSize)
	if r.Header.Get("Create-Directory") != "True" {
		err := createFile(filePath, r)
		if err != nil {
			w.Header().Set("Server", h.cfg.Domain)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Header().Set("Server", h.cfg.Domain)
		w.WriteHeader(http.StatusCreated)
		return

	} else {
		err := createDir(filePath, r.Body)
		if err != nil {
			w.Header().Set("Server", h.cfg.Domain)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		w.Header().Set("Server", h.cfg.Domain)
		w.WriteHeader(http.StatusCreated)
	}
}

func createFile(filePath string, r *http.Request) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return ErrFileFound
	}

	path := strings.Split(filePath, "/")
	_, err = os.Stat(strings.Join(path[:len(path)-1], "/"))
	if err != nil {
		return ErrFileNotFound
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

func createDir(filePath string, r io.Reader) error {
	_, err := os.Stat(filePath)
	if err == nil {
		return ErrFileFound
	}

	path := strings.Split(filePath, "/")
	_, err = os.Stat(strings.Join(path[:len(path)-1], "/"))
	if err != nil {
		return ErrFileNotFound
	}

	err = os.Mkdir(filePath, 0777)
	if err != nil {
		return fmt.Errorf("not exist on disk")
	}
	return nil

}
