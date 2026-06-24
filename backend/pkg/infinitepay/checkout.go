package infinitepay

func (c *Client) CreateCheckout(req CreateCheckoutRequest) (*CreateCheckoutResponse, error) {
	if req.Handle == "" {
		req.Handle = c.handle
	}
	if req.RedirectURL == "" {
		req.RedirectURL = c.redirectURL
	}
	if req.WebhookURL == "" {
		req.WebhookURL = c.webhookURL
	}

	var result CreateCheckoutResponse
	if err := c.do("POST", "/links", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) CheckPayment(req PaymentCheckRequest) (*PaymentCheckResponse, error) {
	if req.Handle == "" {
		req.Handle = c.handle
	}

	var result PaymentCheckResponse
	if err := c.do("POST", "/payment_check", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
