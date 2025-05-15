package domain

import (
	"log/slog"

	"github.com/google/uuid"
)

type UserID uuid.UUID

func (u UserID) LogValue() slog.Value {
	return slog.StringValue(uuid.UUID(u).String())
}

type Pagination struct {
	Offset int
	Limit  int
}

func (p Pagination) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Int("offset", p.Offset),
		slog.Int("limit", p.Limit),
	)
}
