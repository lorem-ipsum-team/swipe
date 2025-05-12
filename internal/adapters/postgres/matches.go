package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/lorem-ipsum-team/swipe/internal/adapters/postgres/gen"
	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

func (r *Repo) GetMatches(ctx context.Context, init domain.UserID, offset, limit int) ([]domain.Match, error) {
	q := gen.New(r.pool)

	dtoMatches, err := q.Matches(ctx,
		gen.MatchesParams{
			InitiatorID: uuid.UUID(init),
			Limit:       int32(limit),   //nolint:gosec
			Offset:      int32(offset)}) //nolint:gosec
	if err != nil {
		return nil, err
	}

	matches := make([]domain.Match, 0, len(dtoMatches))

	for _, dto := range dtoMatches {
		match := domain.Match{
			Init:   domain.UserID(dto.InitiatorID),
			Target: domain.UserID(dto.TargetID),
		}
		matches = append(matches, match)
	}

	return matches, nil
}
