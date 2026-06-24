package abacatepay

func (c *Client) CreateCheckout(req CreateCheckoutRequest) (*Checkout, error) {
	var result Checkout
	if err := c.do("POST", "/checkouts/create", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetCheckout(id string) (*Checkout, error) {
	var result Checkout
	if err := c.do("GET", "/checkouts/get?id="+id, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
