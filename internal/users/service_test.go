package users

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type repositoryMock struct {
	CreateFunc func(ctx context.Context, user User) (User, error)
}

func (r *repositoryMock) Create(ctx context.Context, user User) (User, error) {
	return r.CreateFunc(ctx, user)
}

func (r *repositoryMock) Update(ctx context.Context, email, password string) error {
	return nil
}

func (r *repositoryMock) GetByEMail(ctx context.Context, email string) (User, error) {
	return User{}, nil
}

func (r *repositoryMock) GetAll(ctx context.Context, limit, offset int) ([]User, error) {
	return nil, nil
}

type zipCodeServiceMock struct {
	GetAddressByZipCodeFunc func(zipCode string) (Address, error)
}

func (z *zipCodeServiceMock) GetAddressByZipCode(zipCode string) (Address, error) {
	return z.GetAddressByZipCodeFunc(zipCode)
}

type hashServiceMock struct {
	CreateFunc  func(password string) (string, error)
	IsValidFunc func(inputPassword, password string) bool
}

func (h *hashServiceMock) Create(password string) (string, error) {
	return h.CreateFunc(password)
}

func (h *hashServiceMock) IsValid(inputPassword, password string) bool {
	return h.IsValidFunc(inputPassword, password)
}

func TestService_Create(t *testing.T) {
	tests := []struct {
		name          string
		input         CreateUser
		setupMocks    func(*repositoryMock, *zipCodeServiceMock, *hashServiceMock)
		expectedError error
		expectedUser  User
	}{
		{
			name: "Password mismatch",
			input: CreateUser{
				Email:           "test@example.com",
				Password:        "123456",
				ConfirmPassword: "654321",
				ZipCode:         "12345",
			},
			setupMocks:    func(r *repositoryMock, z *zipCodeServiceMock, h *hashServiceMock) {},
			expectedError: ErrPasswordMismatch,
		},
		{
			name: "Fail to fetch address",
			input: CreateUser{
				Email:           "test@example.com",
				Password:        "123456",
				ConfirmPassword: "123456",
				ZipCode:         "12345",
			},
			setupMocks: func(r *repositoryMock, z *zipCodeServiceMock, h *hashServiceMock) {
				z.GetAddressByZipCodeFunc = func(zipCode string) (Address, error) {
					return Address{}, errors.New("address not found")
				}
			},
			expectedError: errors.New("failed to fetch address: address not found"),
		},
		{
			name: "Fail to hash password",
			input: CreateUser{
				Email:           "test@example.com",
				Password:        "123456",
				ConfirmPassword: "123456",
				ZipCode:         "12345",
			},
			setupMocks: func(r *repositoryMock, z *zipCodeServiceMock, h *hashServiceMock) {
				z.GetAddressByZipCodeFunc = func(zipCode string) (Address, error) {
					return Address{Street: "Main St"}, nil
				}
				h.CreateFunc = func(password string) (string, error) {
					return "", errors.New("hash error")
				}
			},
			expectedError: errors.New("failed to hash password: hash error"),
		},
		{
			name: "Success creating user",
			input: CreateUser{
				Email:           "test@example.com",
				Password:        "123456",
				ConfirmPassword: "123456",
				ZipCode:         "12345",
			},
			setupMocks: func(r *repositoryMock, z *zipCodeServiceMock, h *hashServiceMock) {
				z.GetAddressByZipCodeFunc = func(zipCode string) (Address, error) {
					return Address{Street: "Main St", ZipCode: zipCode}, nil
				}
				h.CreateFunc = func(password string) (string, error) {
					return "hashedPassword", nil
				}
				r.CreateFunc = func(ctx context.Context, user User) (User, error) {
					return user, nil
				}
			},
			expectedError: nil,
			expectedUser: User{
				Email: "test@example.com",
				Address: Address{
					Street:  "Main St",
					ZipCode: "12345",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := &repositoryMock{}
			zipCodeMock := &zipCodeServiceMock{}
			hashMock := &hashServiceMock{}

			tt.setupMocks(repoMock, zipCodeMock, hashMock)

			service := NewService(repoMock, nil, zipCodeMock, hashMock, nil)

			user, err := service.Create(context.Background(), tt.input)

			if tt.expectedError != nil {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError.Error())
				return
			}

			require.Equal(t, tt.expectedUser, user)
			require.NoError(t, err)
		})
	}
}
