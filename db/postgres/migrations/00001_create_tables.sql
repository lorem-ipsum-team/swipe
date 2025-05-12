-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "swipe_db" (
  "initiator_id" uuid NOT NULL,
  "target_id" uuid NOT NULL,
  "initiator_resp" boolean,
  "target_resp" boolean,
  "created_at" timestamp  NOT NULL DEFAULT (now()),
  PRIMARY KEY ("initiator_id", "target_id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "swipe_db";
-- +goose StatementEnd
