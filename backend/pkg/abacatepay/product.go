package abacatepay

func (c *Client) CreateProduct(req CreateProductRequest) (*Product, error) {
	var result Product
	if err := c.do("POST", "/products/create", req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
