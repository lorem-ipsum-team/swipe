package domain

import (
	"log/slog"

	"github.com/lorem-ipsum-team/swipe/pkg/logger"
)

type Swipe struct {
	Init       UserID
	Target     UserID
	InitResp   *bool
	TargetResp *bool
}

func (s Swipe) LogValue() slog.Value {
	return slog.GroupValue(
		logger.NewAttr("init", s.Init.LogValue()),
		logger.NewAttr("target", s.Target.LogValue()),
		logger.NewAttr("initResp", logger.Nullable(s.InitResp)),
		logger.NewAttr("targetResp", logger.Nullable(s.TargetResp)),
	)
}
