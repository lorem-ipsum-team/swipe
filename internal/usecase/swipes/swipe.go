package swipes

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

type SwipeRepo interface {
	CreateSwipe(ctx context.Context, swipe domain.Swipe) error
	MySwipes(ctx context.Context, init domain.UserID, pag domain.Pagination) ([]domain.Swipe, error)
}

type SwipePublisher interface {
	PublishSwipe(ctx context.Context, swipe domain.Swipe) error
}
type UseCase struct {
	repo SwipeRepo
	pub  SwipePublisher
	log  *slog.Logger
}

func New(log *slog.Logger, repo SwipeRepo, pub SwipePublisher) UseCase {
	return UseCase{
		repo: repo,
		pub:  pub,
		log:  log.WithGroup("swipes_usecase"),
	}
}

func (u UseCase) CreateSwipe(ctx context.Context, swipe domain.Swipe) error {
	u.log.DebugContext(ctx, "creating swipe", slog.Any("swipe", swipe))

	err := u.pub.PublishSwipe(ctx, swipe)
	if err != nil {
		return fmt.Errorf("failed to publish swipe (%+#v): %w", swipe, err)
	}

	err = u.repo.CreateSwipe(ctx, swipe)
	if err != nil {
		return fmt.Errorf("failed to create swipe (%+#v): %w", swipe, err)
	}

	return nil
}

func (u UseCase) MySwipes(ctx context.Context, id domain.UserID, pag domain.Pagination) ([]domain.Swipe, error) {
	u.log.DebugContext(ctx, "getting swipes", slog.Any("userID", id), slog.Any("pag", pag))

	swipes, err := u.repo.MySwipes(ctx, id, pag)
	if err != nil {
		return nil, fmt.Errorf("failed to get mySwipes for userID (%s): %w", id, err)
	}

	return swipes, nil
}
