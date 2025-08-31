package cli

import "context"

func (c *Client) CreateAccount(ctx context.Context, accountName, password string) error {
	return c.auth.CreateAccount(ctx, accountName, password)
}

func (c *Client) Login(ctx context.Context, accountName, password string) (err error) {
	if err := c.localStorage.CreateScheme(ctx); err != nil {
		return err
	}

	if c.credential, err = c.auth.Login(ctx, accountName, password, c.publicKeyContent); err != nil {
		return err
	}
	return c.saveCredential(accountName)
}
