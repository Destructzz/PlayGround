package domain

import (
	"backend/internal/repo/sqlc"
)

type CreatePaymentRequest struct {
	BookingID     int64              `json:"booking_id" binding:"required"`
	Amount        string             `json:"amount" binding:"required"`
	Currency      string             `json:"currency" binding:"required"`
	PaymentMethod sqlc.PaymentMethod `json:"payment_method" binding:"required,oneof=cash card online"`
	Status        sqlc.PaymentStatus `json:"status" binding:"omitempty,oneof=pending paid failed refunded"`
	ReceiptNumber string             `json:"receipt_number"`
	PaidAt        string             `json:"paid_at"`
}

type PatchPaymentRequest struct {
	BookingID     *int64              `json:"booking_id" binding:"omitempty"`
	Amount        *string             `json:"amount" binding:"omitempty"`
	Currency      *string             `json:"currency" binding:"omitempty"`
	PaymentMethod *sqlc.PaymentMethod `json:"payment_method" binding:"omitempty,oneof=cash card online"`
	Status        *sqlc.PaymentStatus `json:"status" binding:"omitempty,oneof=pending paid failed refunded"`
	ReceiptNumber *string             `json:"receipt_number"`
	PaidAt        *string             `json:"paid_at"`
}
