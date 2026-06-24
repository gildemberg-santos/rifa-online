package abacatepay

type APIResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Success bool        `json:"success"`
}

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ExternalID  string `json:"externalId,omitempty"`
}

type CreateProductRequest struct {
	ExternalID  string `json:"externalId"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price"`
	Currency    string `json:"currency"`
}

type CheckoutItem struct {
	ID       string `json:"id"`
	Quantity int    `json:"quantity"`
}

type Checkout struct {
	ID           string        `json:"id"`
	URL          string        `json:"url"`
	Status       string        `json:"status"`
	Amount       int           `json:"amount"`
	ExternalID   string        `json:"externalId,omitempty"`
	Items        []CheckoutItem `json:"items,omitempty"`
	DevMode      bool          `json:"devMode"`
}

type CreateCheckoutRequest struct {
	Items         []CheckoutItem `json:"items"`
	ExternalID    string         `json:"externalId,omitempty"`
	ReturnURL     string         `json:"returnUrl,omitempty"`
	CompletionURL string         `json:"completionUrl,omitempty"`
	Methods       []string       `json:"methods,omitempty"`
}
