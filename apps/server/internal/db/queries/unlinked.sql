-- Creating an Unlinked Member reuses CreateMember with a NULL line_user_id.

-- name: LinkMember :one
-- Blank Applicant fields keep what the Admin already recorded for the Unlinked Member.
UPDATE members
SET line_user_id = sqlc.arg(line_user_id),
    cohort       = COALESCE(NULLIF(CAST(sqlc.arg(cohort) AS TEXT), ''), cohort),
    real_name    = COALESCE(NULLIF(CAST(sqlc.arg(real_name) AS TEXT), ''), real_name),
    address      = COALESCE(NULLIF(CAST(sqlc.arg(address) AS TEXT), ''), address),
    phone        = COALESCE(NULLIF(CAST(sqlc.arg(phone) AS TEXT), ''), phone)
WHERE id = sqlc.arg(id) AND line_user_id IS NULL AND deleted_at IS NULL
RETURNING *;
