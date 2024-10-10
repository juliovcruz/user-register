package users

import "errors"

var (
	ErrNotFound   = errors.New("user not found")
	ErrBadRequest = errors.New("error bad request")
)

var (
	ErrMailAlreadyExists = errors.New("mail already exists")
	ErrPasswordMismatch  = errors.New("password and confirm password do not match")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidLogin      = errors.New("invalid email or password")
)

type Address struct {
	Street       string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	Number       string `json:"number"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zip_code"`
}

type CreateUser struct {
	Name            string `json:"name" validate:"required,min=3,max=100" example:"User Name"`
	Email           string `json:"email" validate:"required,email" example:"user@example.com"`
	Password        string `json:"password" validate:"required,min=6,max=100" example:"password"`
	ZipCode         string `json:"zip_code" validate:"required,len=8" example:"74360400"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password" example:"password"`
}

type User struct {
	ID       int64   `json:"id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Address  Address `json:"address"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email" example:"user@example.com"`
	Password string `json:"password" validate:"required,min=6,max=100" example:"password"`
}

type UpdatePassword struct {
	Email           string `json:"email" validate:"required,email" example:"user@example.com"`
	Password        string `json:"password" validate:"required,min=6,max=100" example:"password"`
	ConfirmPassword string `json:"confirm_password" validate:"eqfield=Password" example:"password"`
	Code            int    `json:"code" validate:"required" example:"123456"`
}

type ForgotPassword struct {
	Email string `json:"email" validate:"required,email" example:"user@example.com"`
}
