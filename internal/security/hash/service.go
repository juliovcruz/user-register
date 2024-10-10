package hash

import (
	"github.com/juliovcruz/user-register/internal/settings"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	PreviousSecret string
	CurrentSecret  string
}

func NewService(settings settings.Settings) *Service {
	return &Service{
		PreviousSecret: settings.Database.Secrets.Previous,
		CurrentSecret:  settings.Database.Secrets.Current,
	}
}

func (s *Service) Create(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password+s.CurrentSecret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (s *Service) IsValid(inputPassword, password string) bool {
	secrets := []string{s.CurrentSecret, s.PreviousSecret}

	for _, secret := range secrets {
		if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(inputPassword+secret)); err == nil {
			return true
		}
	}

	return false
}
