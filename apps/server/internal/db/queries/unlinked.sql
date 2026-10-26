-- Creating an Unlinked Member reuses CreateMember with a NULL line_user_id.

-- name: GetOrderOwner :one
SELECT m.* FROM members m JOIN orders o ON o.member_id = m.id WHERE o.id = sqlc.arg(order_id);

-- name: CountMemberOrders :one
SELECT COUNT(*) FROM orders WHERE member_id = sqlc.arg(member_id);

-- name: MoveOrder :execrows
UPDATE orders SET member_id = sqlc.arg(to_member_id) WHERE id = sqlc.arg(order_id);

-- name: MoveMemberOrders :execrows
UPDATE orders SET member_id = sqlc.arg(to_member_id) WHERE member_id = sqlc.arg(from_member_id);

-- name: MoveSponsorPositions :exec
-- Runs before the Unlinked Member is deleted so ON DELETE SET NULL does not drop the assignment.
UPDATE sponsor_positions SET member_id = sqlc.arg(to_member_id) WHERE member_id = sqlc.arg(from_member_id);

-- name: DeleteUnlinkedMember :exec
DELETE FROM members WHERE id = sqlc.arg(id) AND line_user_id IS NULL;
