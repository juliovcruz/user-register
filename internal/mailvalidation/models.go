package mailvalidation

import (
	"errors"
	"time"
)

var (
	ErrRecordNotFound  = errors.New("record not found")
	ErrInvalidCode     = errors.New("invalid code")
	ErrCodeAlreadySent = errors.New("code already sent, wait to send again")
	ErrCodeExpired     = errors.New("code expired, try again")
)

type MailValidation struct {
	Email     string    `json:"email"`
	Code      int       `json:"code"`
	ExpiredAt time.Time `json:"expired_at"`
}
