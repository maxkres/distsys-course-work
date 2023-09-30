package handler

import (
	"compress/gzip"
	"errors"
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func (h Handler) get(w http.ResponseWriter, r *http.Request) {
	filePath := h.cfg.WorkingDirectory + r.URL.Path

	if r.Header.Get("Accept-Encoding") != "gzip" {
		dirInfo, err := getFile(filePath)
		if err != nil {
			if errors.Is(err, ErrFileNotFound) {
				w.Header().Set("Server", h.cfg.Domain)
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Server", h.cfg.Domain)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if dirInfo == nil {
			contentType, err := mimetype.DetectFile(filePath)
			if err != nil {
				w.Header().Set("Server", h.cfg.Domain)
				http.Error(w, "mimetype detect failed", http.StatusNotFound)
				return
			}
			(w).Header().Set("Content-Type", contentType.String())
			w.Header().Set("Server", h.cfg.Domain)
			http.ServeFile(w, r, filePath)
			return
		}

		(w).Header().Set("Content-Type", "text/plain")
		w.Header().Set("Server", h.cfg.Domain)
		_, err = (w).Write(dirInfo)
		if err != nil {
			w.Header().Set("Server", h.cfg.Domain)
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return

	} else {
		(w).Header().Set("Content-Encoding", "gzip")
		err := h.getFileGzip(filePath, w)
		if err != nil {
			log.Println(err)
			w.Header().Set("Server", h.cfg.Domain)
			http.Error(w, fmt.Sprintf("not exist on disk"), http.StatusNotFound)
			return

		}
	}
}

func getFile(filePath string) ([]byte, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, ErrFileNotFound
	}

	if !fileInfo.IsDir() {
		return nil, nil
	}
	cmd := exec.Command("ls", "-lA")

	cmd.Dir = filePath

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("unable to execute command")
	}
	return output, nil
}

func (h Handler) getFileGzip(filePath string, w http.ResponseWriter) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return ErrFileNotFound
	}
	if fileInfo.IsDir() {
		cmd := exec.Command("ls", "-lA")

		cmd.Dir = filePath

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(err)
			return fmt.Errorf("unable to execute command")
		}

		w.Header().Set("Server", h.cfg.Domain)

		gz := gzip.NewWriter(w)
		defer gz.Close()

		if _, err := gz.Write(output); err != nil {
			log.Println(err)
			return fmt.Errorf("unable to compress file: %w", err)
		}

		return nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return ErrFileNotFound
	}

	defer file.Close()

	contentType, err := mimetype.DetectFile(filePath)
	if err != nil {
		return fmt.Errorf("mimetype detect failed")
	}
	w.Header().Set("Content-Type", contentType.String())
	w.Header().Set("Server", h.cfg.Domain)

	gz := gzip.NewWriter(w)
	defer gz.Close()
	if _, err := io.Copy(gz, file); err != nil {
		log.Println(err)
		return fmt.Errorf("unable to compress file: %w", err)
	}

	return nil
}
