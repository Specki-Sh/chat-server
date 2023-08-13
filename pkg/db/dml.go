package db

// User queries
const (
	InsertUserQuery                   = `INSERT INTO users(username, password, email) VALUES ($1, $2, $3) returning id`
	SelectUserByEmailAndPasswordQuery = `SELECT id, username, password, email FROM users WHERE email = $1 AND password = $2`
	SelectUserByIDQuery               = `SELECT id, username, password, email FROM users WHERE id = $1`
	UpdateUserQuery                   = `UPDATE users SET username = $1, password = $2, email = $3 WHERE id = $4`
)

// Member queries
const (
	InsertMemberQuery             = ` INSERT INTO members (room_id, user_id) VALUES ($1, $2)`
	SelectMemberBulkByRoomIDQuery = `SELECT room_id, user_id FROM members WHERE room_id = $1`
	UpdateMemberQuery             = `UPDATE members SET room_id = $1, user_id = $2 WHERE room_id = $3 AND user_id = $4`
	DeleteMemberQuery             = `DELETE FROM members WHERE room_id = $1 AND user_id = $2`
)

// Message queries
const (
	InsertMessageQuery                    = `INSERT INTO messages (sender_id, room_id, content) VALUES ($1, $2, $3) RETURNING id, created_at`
	SelectMessageQuery                    = `SELECT id, sender_id, room_id, content, status, created_at, updated_at, deleted_at FROM messages WHERE id = $1 AND is_active = true`
	UpdateMessageQuery                    = `UPDATE messages SET sender_id = $1, room_id = $2, content = $3, status = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`
	SoftDeleteMessageByIDQuery            = `UPDATE messages SET is_active = false, deleted_at = CURRENT_TIMESTAMP WHERE id = $1`
	SoftDeleteMessageBulkByRoomIDQuery    = `UPDATE messages SET is_active = false, deleted_at = CURRENT_TIMESTAMP WHERE room_id = $1`
	SelectMessageBulkPaginateQuery        = `SELECT id, sender_id, room_id, content, status, created_at, updated_at, deleted_at FROM messages WHERE is_active = true AND room_id = $1 LIMIT $2 OFFSET $3`
	SelectMessageBulkPaginateReverseQuery = `SELECT id, sender_id, room_id, content, status, created_at, updated_at, deleted_at FROM messages WHERE is_active = true AND room_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
)

// Room queries
const (
	InsertRoomQuery     = `INSERT INTO rooms (owner_id, name) VALUES ($1, $2) RETURNING id`
	SelectRoomByIDQuery = `SELECT id, owner_id, name FROM rooms WHERE id = $1`
	UpdateRoomQuery     = `UPDATE rooms SET name = $1 WHERE id = $2`
	DeleteRoomQuery     = `DELETE FROM rooms WHERE id = $1`
)
