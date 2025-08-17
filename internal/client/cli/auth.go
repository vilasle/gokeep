package cli

import "context"

func (c *Client) CreateAccount(ctx context.Context, accountName, password string) error {
	if err := c.auth.CreateAccount(accountName, password, c.publicKeyContent); err != nil {
		return err
	}

	return c.localStorage.CreateScheme(ctx)
}

func (c *Client) Login(accountName, password string) (err error) {
	if len(c.credential) > 0 {
		return nil
	}

	c.credential, err = c.auth.Login(accountName, password)

	if err := c.saveCredential(accountName); err != nil {
		return err
	}

	return err
}
