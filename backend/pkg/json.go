package pkg

import (
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
)

func JSONOrObject(b []byte) json.RawMessage {
	if len(b) == 0 {
		return json.RawMessage("{}")
	}

	return json.RawMessage(b)
}

func JSONOrArray(b []byte) json.RawMessage {
	if len(b) == 0 {
		return json.RawMessage("[]")
	}

	return json.RawMessage(b)
}

func TextValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}

	return value.String
}

func NumericToString(value pgtype.Numeric) string {
	if !value.Valid {
		return ""
	}

	raw, err := value.Value()
	if err != nil {
		return ""
	}

	switch v := raw.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}
