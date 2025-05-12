
-- name: Matches :many
SELECT 
    m.initiator_id,
    m.target_id
FROM 
    user_matches m
WHERE
    m.initiator_id = $1 OR m.target_id = $1
ORDER BY m.initiator_id DESC LIMIT $2 OFFSET $3;
