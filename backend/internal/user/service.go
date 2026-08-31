package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ErrEmailTaken = errors.New("Emails already registered")
var ErrInvalidCredentials = errors.New("invalid email or password")
var ErrQuotaExceeded = errors.New("daily search limit reached")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(email, password string) (*User, error) {
	_, err := s.repo.FindByEmail(email)
	if err == nil {
		return nil, ErrEmailTaken
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	u := &User{Email: email, PasswordHash: string(hash)}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) Login(email, password string) (*User, error) {
	u, err := s.repo.FindByEmail(email)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return u, nil
}

func (s *Service) GetByID(id string) (*User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) ConsumeSearch(userID string) error {
	return s.repo.ConsumeSearch(userID)
}
