package matches

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

type matchesRepo interface {
	GetMatches(ctx context.Context, init domain.UserID, pag domain.Pagination) ([]domain.Match, error)
}
type UseCase struct {
	repo matchesRepo
	log  *slog.Logger
}

func New(log *slog.Logger, repo matchesRepo) UseCase {
	return UseCase{
		repo: repo,
		log:  log.WithGroup("matches_usecase"),
	}
}

func (u UseCase) Matches(ctx context.Context, id domain.UserID, pag domain.Pagination) ([]domain.Match, error) {
	u.log.DebugContext(ctx, "getting matches", slog.Any("userID", id), slog.Any("pag", pag))

	matches, err := u.repo.GetMatches(ctx, id, pag)
	if err != nil {
		return nil, fmt.Errorf("failed to get matches for userID (%s): %w", id, err)
	}

	return matches, nil
}
