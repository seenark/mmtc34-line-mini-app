-- name: GetMemberByLineUserID :one
SELECT * FROM members WHERE line_user_id = sqlc.arg(line_user_id);

-- name: GetMember :one
SELECT * FROM members WHERE id = sqlc.arg(id);

-- name: ListMembers :many
SELECT * FROM members
WHERE CAST(sqlc.arg(status) AS TEXT) = '' OR status = sqlc.arg(status)
ORDER BY id;

-- name: CreateMember :one
INSERT INTO members (line_user_id, status, role, cohort, nickname, real_name, address, phone, created_at)
VALUES (sqlc.arg(line_user_id), sqlc.arg(status), sqlc.arg(role), sqlc.arg(cohort), sqlc.arg(nickname),
        sqlc.arg(real_name), sqlc.arg(address), sqlc.arg(phone), sqlc.arg(now))
RETURNING *;

-- name: UpsertSuperAdmin :exec
INSERT INTO members (line_user_id, status, role, nickname, created_at)
VALUES (sqlc.arg(line_user_id), 'approved', 'admin', sqlc.arg(nickname), sqlc.arg(now))
ON CONFLICT (line_user_id) DO UPDATE SET status = 'approved', role = 'admin';

-- name: ApproveMember :one
UPDATE members SET status = 'approved' WHERE id = sqlc.arg(id) AND deleted_at IS NULL RETURNING *;

-- name: DeleteApplicant :exec
DELETE FROM members WHERE id = sqlc.arg(id) AND status = 'applicant';
