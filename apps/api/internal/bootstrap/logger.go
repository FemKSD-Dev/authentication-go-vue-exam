package bootstrap

import (
	"github.com/thessem/zap-prettyconsole"
	"go.uber.org/zap"
)

func NewLogger() (*zap.Logger, error) {
	return prettyconsole.NewLogger(zap.DebugLevel), nil
}
