package zipcode

import "github.com/juliovcruz/user-register/internal/users"

type Client interface {
	GetAddressByZipCode(zipCode string) (users.Address, error)
}

type Service struct {
	client Client
}

func NewService(client Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) GetAddressByZipCode(zipCode string) (users.Address, error) {
	return s.client.GetAddressByZipCode(zipCode)
}
