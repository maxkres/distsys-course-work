package service

import (
	"auth/config"
	"auth/internal/model"
	"crypto/sha1"
	"fmt"
	"os"
)

type Service struct {
	r   Repo
	cfg *config.Config
}

type Repo interface {
	SignUp(login, password string) error
	GetPassword(login string) (string, error)
}

func New(redis Repo, cfg *config.Config) *Service {
	return &Service{
		redis,
		cfg,
	}
}

func (s *Service) SignUp(user model.User) (string, error) {
	var err error
	password, err := s.GenerateHash(user.Username, user.Password)
	if err != nil {
		return "", fmt.Errorf("generate hash failed: %w", err)
	}

	err = s.r.SignUp(user.Username, password)
	if err != nil {
		return "", err
	}

	prvKey, err := os.ReadFile(s.cfg.PemPath)
	if err != nil {
		return "", err
	}
	token := NewJWT(prvKey, nil)
	if err != nil {
		return "", fmt.Errorf("new token failed: %w", err)
	}

	tok, err := token.Create(user.Username)
	if err != nil {
		return "", err
	}

	return tok, nil
}

func (s *Service) GenerateHash(username, password string) (string, error) {
	hash := sha1.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}
	return string(hash.Sum([]byte(username))), nil
}

func (s *Service) SingIn(user model.User) (string, error) {
	passwordDB, err := s.r.GetPassword(user.Username)
	if err != nil {
		return "", fmt.Errorf("check user by phone number failed: %w", err)
	}

	hash := sha1.New()
	_, err = hash.Write([]byte(user.Password))
	if err != nil {
		return "", fmt.Errorf("write failed: %w", err)
	}

	if passwordDB != string(hash.Sum([]byte(user.Username))) {
		return "", fmt.Errorf("incorrect password")
	}

	prvKey, err := os.ReadFile(s.cfg.PemPath)
	if err != nil {
		return "", err
	}

	token := NewJWT(prvKey, nil)
	if err != nil {
		return "", fmt.Errorf("new token failed: %w", err)
	}

	tok, err := token.Create(user.Username)
	if err != nil {
		return "", err
	}

	return tok, nil
}

func (s *Service) Check(user model.User) error {
	_, err := s.r.GetPassword(user.Username)
	if err == nil {
		return fmt.Errorf("user already exists")
	}

	return nil
}

func (s *Service) Verify(tok string) (string, error) {
	pubKey, err := os.ReadFile(s.cfg.PubPath)
	if err != nil {
		return "", err
	}

	token := NewJWT(nil, pubKey)
	if err != nil {
		return "", fmt.Errorf("new token failed: %w", err)
	}

	dat, err := token.Validate(tok)
	if err != nil {
		return "", err
	}

	return dat.(string), nil
}
