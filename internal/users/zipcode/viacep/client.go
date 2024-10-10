package viacep

import (
	"encoding/json"
	"fmt"

	"github.com/juliovcruz/user-register/internal/settings"
	"github.com/juliovcruz/user-register/internal/users"
	"github.com/juliovcruz/user-register/internal/users/zipcode"
	"github.com/valyala/fasthttp"
)

type Client struct {
	baseURL    string
	httpClient *fasthttp.Client
}

func NewClient(settings settings.ZipCode) *Client {
	return &Client{
		baseURL:    settings.ViaCEPBaseURL,
		httpClient: &fasthttp.Client{},
	}
}

func (c *Client) GetAddressByZipCode(zipCode string) (users.Address, error) {
	resp := fasthttp.AcquireResponse()
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseResponse(resp)
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s/%s/json", c.baseURL, zipCode))
	if err := c.httpClient.Do(req, resp); err != nil {
		return users.Address{}, fmt.Errorf("request error: %w", err)
	}

	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		{
			var apiResp apiResponse
			if err := json.Unmarshal(resp.Body(), &apiResp); err != nil {
				return users.Address{}, fmt.Errorf("json parse error: %w", err)
			}

			if apiResp.Erro == "true" {
				return users.Address{}, zipcode.ErrZipCodeNotFound
			}

			return parseResponse(apiResp), nil
		}
	case fasthttp.StatusBadRequest:
		return users.Address{}, zipcode.ErrInvalidZipCode
	default:
		return users.Address{}, fmt.Errorf("invalid status_code: %d", resp.StatusCode())
	}
}

func parseResponse(apiResp apiResponse) users.Address {
	return users.Address{
		Street:       apiResp.Logradouro,
		Neighborhood: apiResp.Bairro,
		City:         apiResp.Localidade,
		State:        apiResp.Uf,
		ZipCode:      apiResp.Cep,
	}
}
