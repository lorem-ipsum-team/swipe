package matches

import (
	"context"
	"fmt"

	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

type matchesRepo interface {
	GetMatches(ctx context.Context, init domain.UserID, offset, limit int) ([]domain.Match, error)
}
type UseCase struct {
	repo matchesRepo
}

func (u *UseCase) Matches(ctx context.Context, id domain.UserID, offset, limit int) ([]domain.Match, error) {
	matches, err := u.repo.GetMatches(ctx, id, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches for userID (%s): %w", id, err)
	}

	return matches, nil
}
