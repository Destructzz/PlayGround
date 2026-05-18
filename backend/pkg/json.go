package pkg

import (
	"encoding/json"
	"strconv"

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

func NumericToFloat64(value pgtype.Numeric) float64 {
	if !value.Valid {
		return 0
	}

	s := NumericToString(value)
	if s == "" {
		return 0
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return f
}

