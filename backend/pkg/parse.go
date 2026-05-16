package pkg

import (
	"strconv"
	"time"
)

func ParsePositiveInt64(raw string) (int64, error) {
	value, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func ParseDateYYYYMMDD(raw string) (time.Time, error) {
	return time.Parse("2006-01-02", raw)
}
