package dto

type CustomerDetail struct {
	GivenNames string `json:"given_names"`
	Email      string `json:"email"`
}

type SendInvoice struct {
	ExternalID     string  `json:"external_id"`
	Amount         float32 `json:"amount"`
	Description    string  `json:"description"`
	CustomerDetail `json:"customer"`
}

type SendInvoiceResponse struct {
	InvoiceId  string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Status     string  `json:"status"`
	Amount     float32 `json:"amount"`
	InvoiceURL string  `json:"invoice_url"`
}

type SendInvoiceResponseLessDetailed struct {
	InvoiceId  string  `json:"id"`
	Amount     float32 `json:"amount"`
	InvoiceURL string  `json:"invoice_url"`
	Status     string  `json:"status,omitempty"`
}

type XenditWebhook struct {
	ExternalId    string `json:"external_id"`
	InvoiceId     string `json:"id"`
	Status        string `json:"status"`
	PaymentMethod string `json:"payment_method"`
	CompletedAt   string `json:"completed_at"`
}
