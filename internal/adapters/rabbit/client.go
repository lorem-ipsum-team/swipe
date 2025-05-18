package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/lorem-ipsum-team/swipe/internal/domain"
	"github.com/lorem-ipsum-team/swipe/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	msgTimeout = 5 * time.Second
)

type Repo struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	log        *slog.Logger
	swipeQueue string
}

type Swipe struct {
	Init   uuid.UUID `json:"init"`
	Target uuid.UUID `json:"target"`
	Like   bool      `json:"like"`
}

func New(ctx context.Context, log *slog.Logger, url, swipeQueue string) (*Repo, error) {
	log = log.WithGroup("rabbit_repo")
	log.Info("connect to rabbit", slog.Any("connection string", logger.Secret(url)))

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Rabbit: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		swipeQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	repo := &Repo{
		conn:       conn,
		channel:    channel,
		log:        log,
		swipeQueue: swipeQueue,
	}

	go func() {
		<-ctx.Done()

		_ = repo.Close()
	}()

	return repo, nil
}

func (r *Repo) Close() error {
	r.log.Info("closing rabbit_repo")

	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close amqp channel: %w", err)
	}

	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close amqp connection: %w", err)
	}

	return nil
}

func (r *Repo) PublishSwipe(ctx context.Context, swipe domain.Swipe) error {
	swipeDto := Swipe{
		Init:   uuid.UUID(swipe.Init),
		Target: uuid.UUID(swipe.Target),
		Like:   *swipe.InitResp,
	}

	body, err := json.Marshal(swipeDto)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, msgTimeout)
	defer cancel()

	err = r.channel.PublishWithContext(
		ctx,
		"",
		r.swipeQueue,
		false,
		false,
		amqp.Publishing{ //nolint:exhaustruct
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish swipe: %w", err)
	}

	return nil
}
