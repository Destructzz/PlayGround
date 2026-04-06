package service

import (
	"context"
	"time"

	"backend/internal/domain"
	"backend/internal/repo/sqlc"

	"github.com/jackc/pgx/v5/pgtype"
)

type PaymentService struct {
	queries *sqlc.Queries
}

func NewPayment(queries *sqlc.Queries) *PaymentService {
	return &PaymentService{queries: queries}
}

func (p *PaymentService) CreatePayment(ctx context.Context, dto domain.CreatePaymentRequest) (sqlc.Payment, error) {
	var parsedAmount pgtype.Numeric
	if err := parsedAmount.Scan(dto.Amount); err != nil {
		return sqlc.Payment{}, err
	}

	params := sqlc.CreatePaymentParams{
		BookingID:     dto.BookingID,
		Amount:        parsedAmount,
		Currency:      dto.Currency,
		PaymentMethod: dto.PaymentMethod,
		Status:        dto.Status,
		ReceiptNumber: pgtype.Text{
			String: dto.ReceiptNumber,
			Valid:  dto.ReceiptNumber != "",
		},
	}

	if dto.PaidAt != "" {
		t, err := time.Parse(time.RFC3339, dto.PaidAt)
		if err != nil {
			return sqlc.Payment{}, err
		}
		params.PaidAt = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	return p.queries.CreatePayment(ctx, params)
}

func (p *PaymentService) ListPayments(ctx context.Context) ([]sqlc.Payment, error) {
	return p.queries.ListPayments(ctx)
}

func (p *PaymentService) GetPaymentByID(ctx context.Context, id int64) (sqlc.Payment, error) {
	return p.queries.GetPaymentByID(ctx, id)
}

func (p *PaymentService) PatchPayment(ctx context.Context, id int64, dto domain.PatchPaymentRequest) (sqlc.Payment, error) {
	params := sqlc.PatchPaymentParams{ID: id}

	if dto.BookingID != nil {
		params.BookingID = pgtype.Int8{
			Int64: *dto.BookingID,
			Valid: true,
		}
	}

	if dto.Amount != nil {
		var parsedAmount pgtype.Numeric
		if err := parsedAmount.Scan(*dto.Amount); err != nil {
			return sqlc.Payment{}, err
		}
		params.Amount = parsedAmount
	}

	if dto.Currency != nil {
		params.Currency = pgtype.Text{
			String: *dto.Currency,
			Valid:  true,
		}
	}

	if dto.PaymentMethod != nil {
		params.PaymentMethod = sqlc.NullPaymentMethod{
			PaymentMethod: *dto.PaymentMethod,
			Valid:         true,
		}
	}

	if dto.Status != nil {
		params.Status = sqlc.NullPaymentStatus{
			PaymentStatus: *dto.Status,
			Valid:         true,
		}
	}

	if dto.ReceiptNumber != nil {
		params.ReceiptNumber = pgtype.Text{
			String: *dto.ReceiptNumber,
			Valid:  true,
		}
	}

	if dto.PaidAt != nil {
		t, err := time.Parse(time.RFC3339, *dto.PaidAt)
		if err != nil {
			return sqlc.Payment{}, err
		}
		params.PaidAt = pgtype.Timestamptz{
			Time:  t,
			Valid: true,
		}
	}

	return p.queries.PatchPayment(ctx, params)
}

func (p *PaymentService) DeletePayment(ctx context.Context, id int64) (int64, error) {
	return p.queries.DeletePayment(ctx, id)
}
