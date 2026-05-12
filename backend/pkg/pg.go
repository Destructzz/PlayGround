package pkg

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TextPointer(value pgtype.Text) *string {
	if !value.Valid {
		return nil
	}

	text := value.String
	return &text
}

func UUIDString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}

	return uuid.UUID(value.Bytes).String()
}
