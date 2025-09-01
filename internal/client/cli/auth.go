package cli

import "context"

//CreateAccount - create new account on server
func (c *CommandLineClient) CreateAccount(ctx context.Context, accountName, password string) error {
	return c.auth.CreateAccount(ctx, accountName, password)
}

//Login - login to server and save credential for next requests
func (c *CommandLineClient) Login(ctx context.Context, accountName, password string) (err error) {
	if err := c.localStorage.CreateScheme(ctx); err != nil {
		return err
	}

	if c.credential, err = c.auth.Login(ctx, accountName, password, c.publicKeyContent); err != nil {
		return err
	}
	return c.saveCredential(accountName)
}
