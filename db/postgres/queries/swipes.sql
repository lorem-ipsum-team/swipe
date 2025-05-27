-- name: SwipeExists :one
SELECT EXISTS(
    SELECT 1 FROM swipe_db 
    WHERE initiator_id = $1 AND target_id = $2
);

-- name: UpsertInitSwipe :exec
INSERT INTO swipe_db (initiator_id, target_id, initiator_resp)
VALUES ($1, $2, $3)
ON CONFLICT (initiator_id, target_id) 
DO UPDATE SET 
    initiator_resp = EXCLUDED.initiator_resp;

-- name: UpsertTargetSwipe :exec
INSERT INTO swipe_db (initiator_id, target_id, target_resp)
VALUES ($1, $2, $3)
ON CONFLICT (initiator_id, target_id) 
DO UPDATE SET 
    target_resp = EXCLUDED.target_resp;


-- name: SwipesTargetLike :many
SELECT 
    initiator_id,
    target_id,
    initiator_resp,
    target_resp,
    created_at
FROM swipe_db
WHERE target_id = $1 
    AND initiator_resp = TRUE
    AND target_resp IS NULL
ORDER BY created_at DESC LIMIT $2 OFFSET $3;
