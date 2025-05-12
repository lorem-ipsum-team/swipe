package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/lorem-ipsum-team/swipe/internal/adapters/postgres/gen"
	"github.com/lorem-ipsum-team/swipe/internal/domain"
)

func (r *Repo) CreateSwipe(ctx context.Context, init, target domain.UserID, like bool) (err error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	query := gen.New(tx)

	existsOpposite, err := query.SwipeExists(ctx, gen.SwipeExistsParams{
		InitiatorID: uuid.UUID(target),
		TargetID:    uuid.UUID(init)})
	if err != nil {
		return err
	}

	if existsOpposite {
		err = query.UpsertTargetSwipe(ctx, gen.UpsertTargetSwipeParams{
			InitiatorID: uuid.UUID(target),
			TargetID:    uuid.UUID(init),
			TargetResp:  pgtype.Bool{Bool: like, Valid: true}})
	} else {
		err = query.UpsertInitSwipe(ctx, gen.UpsertInitSwipeParams{
			InitiatorID:   uuid.UUID(init),
			TargetID:      uuid.UUID(target),
			InitiatorResp: pgtype.Bool{Bool: like, Valid: true},
		})
	}

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repo) MySwipes(ctx context.Context, init domain.UserID, offset, limit int) ([]domain.Swipe, error) {
	query := gen.New(r.pool)

	dtoSwipes, err := query.SwipesTargetLike(ctx, gen.SwipesTargetLikeParams{
		TargetID: uuid.UUID(init),
		Limit:    int32(limit),   //nolint:gosec
		Offset:   int32(offset)}) //nolint:gosec
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
