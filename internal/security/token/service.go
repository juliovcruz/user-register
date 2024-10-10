package token

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/juliovcruz/user-register/internal/settings"
	"github.com/juliovcruz/user-register/internal/users"
)

type Service struct {
	Secret string
}

func NewService(settings settings.Settings) *Service {
	return &Service{Secret: settings.TokenSecret}
}

func (s *Service) Create(user users.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Email,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *Service) IsValid(tokenStr string) (bool, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		return false, errors.New("error to parse token")
	}

	return token.Valid, nil
}
