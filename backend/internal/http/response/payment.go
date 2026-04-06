package response

import (
	"time"
)

type DeletePaymentResponse struct {
	ID        int64     `json:"id" example:"1"`
	Timestamp time.Time `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string    `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type PaymentResponse struct {
	Payment   PaymentDoc `json:"payment"`
	Timestamp time.Time  `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string     `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type PaymentListResponse struct {
	Payments  []PaymentDoc `json:"payments"`
	Timestamp time.Time    `json:"timestamp" example:"2026-01-19T15:37:27.514667373Z"`
	RequestID string       `json:"request_id,omitempty" example:"7fbd6854-8e42-4451-80ee-6da60aeceacd"`
}

type PaymentDoc struct {
	ID            int64      `json:"id" example:"1"`
	BookingID     int64      `json:"booking_id" example:"1"`
	Amount        string     `json:"amount" example:"1500.00"`
	Currency      string     `json:"currency" example:"RUB"`
	PaymentMethod string     `json:"payment_method" example:"card"`
	Status        string     `json:"status" example:"pending"`
	ReceiptNumber *string    `json:"receipt_number,omitempty" example:"RCP-001"`
	PaidAt        *time.Time `json:"paid_at,omitempty" example:"2026-03-22T15:00:00Z"`
	CreatedAt     time.Time  `json:"created_at" example:"2026-01-19T15:37:27.514667373Z"`
	UpdatedAt     time.Time  `json:"updated_at" example:"2026-01-19T15:37:27.514667373Z"`
}
