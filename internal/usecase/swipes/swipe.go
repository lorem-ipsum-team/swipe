package swipes

import (
	"context"
	"fmt"

	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

type SwipeRepo interface {
	CreateSwipe(ctx context.Context, init, target domain.UserID, like bool) error
	MySwipes(ctx context.Context, init domain.UserID, offset, limit int) ([]domain.Swipe, error)
}
type UseCase struct {
	repo SwipeRepo
}

func (u *UseCase) CreateSwipe(ctx context.Context, init, target domain.UserID, like bool) error {
	err := u.repo.CreateSwipe(ctx, init, target, like)
	if err != nil {
		return fmt.Errorf("failed to create swipe init (%s) target (%s) like (%t): %w", init, target, like, err)
	}

	return nil
}

func (u *UseCase) MySwipes(ctx context.Context, init domain.UserID, offset, limit int) ([]domain.Swipe, error) {
	swipes, err := u.repo.MySwipes(ctx, init, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get mySwipes for userID (%s): %w", init, err)
	}

	return swipes, nil
}
