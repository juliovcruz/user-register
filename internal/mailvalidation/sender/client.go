package sender

import "context"

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Send(ctx context.Context, email string, code int) error {
	println("sending email to", email, "with code", code)
	return nil
}
