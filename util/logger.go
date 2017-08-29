package util

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// NewDefaultLogger creates a new logger with our prefered options
func NewDefaultLogger() (*zap.Logger, error) {
	config := zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		Encoding:          "json",
		ErrorOutputPaths:  []string{"stdout"},
		Level:             zap.NewAtomicLevel(),
		OutputPaths:       []string{"stdout"},
	}
	l, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create logger")
	}
	return l, nil
}
