package service

import (
	"fmt"
	"store/config"
	"strings"
)

var (
	ErrExists = fmt.Errorf("already exists")
)

type Service struct {
	r   Repo
	cfg *config.Config
}

type Repo interface {
	Store(key, value string) error
	Get(key string) (string, error)
}

func New(redis Repo, cfg *config.Config) *Service {
	return &Service{
		redis,
		cfg,
	}
}

func (s *Service) Store(username, key, value string) error {
	val, _ := s.Get(key)
	if val != "" && !strings.Contains(val, username) {
		return ErrExists
	}
	value = fmt.Sprintf("%s:%s", username, value)
	return s.r.Store(key, value)
}

func (s *Service) Get(key string) (string, error) {
	return s.r.Get(key)
}
