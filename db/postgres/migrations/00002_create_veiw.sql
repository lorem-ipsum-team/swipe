-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE VIEW user_matches AS
SELECT 
    s1.initiator_id AS initiator_id,
    s1.target_id AS target_id
FROM 
    swipe_db s1
WHERE 
    s1.initiator_resp = TRUE AND s1.target_resp = TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS user_matches;
-- +goose StatementEnd
