
-- name: UpdateMemberProfile :one
UPDATE members
SET cohort = sqlc.arg(cohort), nickname = sqlc.arg(nickname), real_name = sqlc.arg(real_name),
    address = sqlc.arg(address), phone = sqlc.arg(phone)
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *;

-- name: SetMemberRole :one
UPDATE members SET role = sqlc.arg(role)
WHERE id = sqlc.arg(id) AND deleted_at IS NULL
RETURNING *;

-- name: EraseMember :exec
-- PDPA erase: the row stays for nickname-keyed Orders; the CHECK on members enforces the blanks.
UPDATE members
SET line_user_id = NULL, cohort = '', real_name = '', address = '', phone = '', deleted_at = sqlc.arg(now)
WHERE id = sqlc.arg(id);

-- name: EraseMemberShipping :exec
UPDATE orders SET ship_name = '', ship_address = '', ship_phone = '' WHERE member_id = sqlc.arg(member_id);
