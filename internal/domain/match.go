package domain

import (
	"log/slog"

	"github.com/lorem-ipsum-team/swipe/pkg/logger"
)

type Match struct {
	Init   UserID
	Target UserID
}

func (m Match) LogValue() slog.Value {
	return slog.GroupValue(
		logger.NewAttr("init", m.Init.LogValue()),
		logger.NewAttr("target", m.Target.LogValue()),
	)
}
