package zipcode

import "errors"

var (
	ErrInvalidZipCode  = errors.New("invalid zip code")
	ErrZipCodeNotFound = errors.New("zip code not found")
)
