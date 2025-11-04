package smartapi

func (c *Client) GenerateSession() error {

	type loginRequest struct {
		ClientCode string `json:"clientcode"`
		Password   string `json:"password"`
		TOTP       string `json:"totp"`
		State      string `json:"state"`
	}

	type loginResponse struct {
		JWTToken     string `json:"jwtToken"`
		RefreshToken string `json:"refreshToken"`
		FeedToken    string `json:"feedToken"`
		State        string `json:"state"`
	}

	var loginReq = loginRequest{
		ClientCode: c.ClientId,
		Password:   c.Mpin,
		TOTP:       c.TOTPKey,
		State:      "1",
	}

	var lr loginResponse

	err := c.doRequest("POST", LOGIN_ENDPOINT, loginReq, &lr)
	if err != nil {
		return err
	}

	c.JWTToken = lr.JWTToken
	c.RefreshToken = lr.RefreshToken
	c.FeedToken = lr.FeedToken

	return nil
}

func (c *Client) UserProfile() (map[string]any, error) {

	return nil, nil
}

func (c *Client) Margins() (map[string]any, error) {
	return nil, nil
}

func (c *Client) Orders() error {

	return nil
}

func (c *Client) Positions() error {
	return nil
}

func (c *Client) TradeBook() error {
	return nil
}
