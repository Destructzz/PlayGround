package observability

import (
	"os"
	"strings"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(env string) (*zap.Logger, error) {
	format := strings.ToLower(strings.TrimSpace(os.Getenv("APP_LOG_FORMAT")))
	level := strings.ToLower(strings.TrimSpace(os.Getenv("APP_LOG_LEVEL")))
	service := strings.TrimSpace(os.Getenv("APP_SERVICE_NAME"))
	if service == "" {
		service = "playground-backend"
	}

	var cfg zap.Config
	switch format {
	case "console":
		cfg = zap.NewDevelopmentConfig()
	case "json":
		cfg = zap.NewProductionConfig()
	default:
		if strings.EqualFold(env, "production") {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
		}
	}

	if level != "" {
		var parsed zapcore.Level
		if err := parsed.Set(level); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(parsed)
		}
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("service", service),
		zap.String("env", strings.ToLower(env)),
	)

	return logger, nil
}
