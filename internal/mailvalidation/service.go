package mailvalidation

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type Repository interface {
	CreateOrUpdate(ctx context.Context, mailValidation MailValidation) error
	GetByEmail(ctx context.Context, email string) (MailValidation, error)
	Delete(ctx context.Context, email string) error
}

type Client interface {
	Send(ctx context.Context, email string, code int) error
}

type Service struct {
	repo      Repository
	expiredIn time.Duration
	client    Client
}

func NewService(repo Repository, client Client) *Service {
	return &Service{repo: repo, expiredIn: 2 * time.Hour, client: client}
}

func (s *Service) Create(ctx context.Context, email string) error {
	alreadyExists, err := s.repo.GetByEmail(ctx, email)
	if !errors.Is(err, ErrRecordNotFound) {
		return err
	}

	if !alreadyExists.ExpiredAt.IsZero() && alreadyExists.ExpiredAt.Before(time.Now()) {
		return ErrCodeAlreadySent
	}

	code := rand.Intn(999999)

	if err := s.client.Send(ctx, email, code); err != nil {
		return err
	}

	return s.repo.CreateOrUpdate(ctx, MailValidation{
		Email:     email,
		Code:      code,
		ExpiredAt: time.Now().Add(s.expiredIn),
	})
}

func (s *Service) Validate(ctx context.Context, email string, code int) error {
	mailValidation, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}
	if mailValidation.Code != code {
		return ErrInvalidCode
	}
	if time.Now().After(mailValidation.ExpiredAt) {
		return ErrCodeExpired
	}

	if err := s.repo.Delete(ctx, email); err != nil {
		return err
	}

	return nil
}
