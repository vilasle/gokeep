package client

func (c *Client) CreateAccount(accountName, password string) error {
	return c.auth.CreateAccount(accountName, password, c.publicKeyContent)
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
