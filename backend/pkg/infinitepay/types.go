package infinitepay

type CreateCheckoutRequest struct {
	Handle      string          `json:"handle"`
	Items       []CheckoutItem  `json:"items"`
	OrderNSU    string          `json:"order_nsu,omitempty"`
	RedirectURL string          `json:"redirect_url,omitempty"`
	WebhookURL  string          `json:"webhook_url,omitempty"`
	Customer    *Customer       `json:"customer,omitempty"`
}

type CheckoutItem struct {
	Quantity    int    `json:"quantity"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

type Customer struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

type CreateCheckoutResponse struct {
	URL string `json:"url"`
}

type PaymentCheckRequest struct {
	Handle         string `json:"handle"`
	OrderNSU       string `json:"order_nsu"`
	TransactionNSU string `json:"transaction_nsu,omitempty"`
	Slug           string `json:"slug,omitempty"`
}

type PaymentCheckResponse struct {
	Success       bool   `json:"success"`
	Paid          bool   `json:"paid"`
	Amount        int    `json:"amount"`
	PaidAmount    int    `json:"paid_amount"`
	Installments  int    `json:"installments"`
	CaptureMethod string `json:"capture_method"`
}

type WebhookPayload struct {
	InvoiceSlug    string           `json:"invoice_slug"`
	Amount         int              `json:"amount"`
	PaidAmount     int              `json:"paid_amount"`
	Installments   int              `json:"installments"`
	CaptureMethod  string           `json:"capture_method"`
	TransactionNSU string           `json:"transaction_nsu"`
	OrderNSU       string           `json:"order_nsu"`
	ReceiptURL     string           `json:"receipt_url"`
	Items          []CheckoutItem   `json:"items"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
