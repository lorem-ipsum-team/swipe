package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lorem-ipsum-team/swipe/internal/adapters/postgres/gen"
	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

func (r *Repo) CreateSwipe(ctx context.Context, swipe domain.Swipe) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	query := gen.New(tx)

	existsOpposite, err := query.SwipeExists(ctx, gen.SwipeExistsParams{
		InitiatorID: uuid.UUID(swipe.Target),
		TargetID:    uuid.UUID(swipe.Init),
	})
	if err != nil {
		return fmt.Errorf("failed to check if swipe exists: %w", err)
	}

	if existsOpposite {
		slog.DebugContext(ctx, "swipe already exists", slog.Any("swipe", swipe))
		err = query.UpsertTargetSwipe(ctx, gen.UpsertTargetSwipeParams{
			InitiatorID: uuid.UUID(swipe.Target),
			TargetID:    uuid.UUID(swipe.Init),
			TargetResp:  pgtype.Bool{Bool: *swipe.InitResp, Valid: true},
		})
	} else {
		err = query.UpsertInitSwipe(ctx, gen.UpsertInitSwipeParams{
			InitiatorID:   uuid.UUID(swipe.Init),
			TargetID:      uuid.UUID(swipe.Target),
			InitiatorResp: pgtype.Bool{Bool: *swipe.InitResp, Valid: true},
		})
	}

	if err != nil {
		return fmt.Errorf("failed to upsert swipe: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *Repo) MySwipes(ctx context.Context, init domain.UserID, pag domain.Pagination) ([]domain.Swipe, error) {
	query := gen.New(r.pool)

	dtoSwipes, err := query.SwipesTargetLike(ctx, gen.SwipesTargetLikeParams{
		TargetID: uuid.UUID(init),
		Limit:    int32(pag.Limit),  //nolint:gosec
		Offset:   int32(pag.Offset), //nolint:gosec
	})
	if err != nil {
		return nil, err
	}

	swipes := make([]domain.Swipe, 0, len(dtoSwipes))

	for _, dto := range dtoSwipes {
		swipe := domain.Swipe{
			Init:       domain.UserID(dto.InitiatorID),
			Target:     domain.UserID(dto.TargetID),
			InitResp:   nil,
			TargetResp: nil,
		}
		if dto.InitiatorResp.Valid {
			swipe.InitResp = &dto.InitiatorResp.Bool
		}

		if dto.TargetResp.Valid {
			swipe.TargetResp = &dto.TargetResp.Bool
		}

		swipes = append(swipes, swipe)
	}

	return swipes, nil
}
