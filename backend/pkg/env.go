package pkg

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

func GetEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}

	zap.L().Warn("missing env var, using fallback", zap.String("key", key), zap.String("fallback", fallback))
	
	return fallback
}