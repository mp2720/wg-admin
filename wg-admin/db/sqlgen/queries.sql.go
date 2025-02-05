// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: queries.sql

package sqlgen

import (
	"context"
	"time"
)

const createUser = `-- name: CreateUser :execlastid


INSERT INTO users (
    uuid,
    name,
    is_admin,
    is_banned,
    fare,
    private_key,
    public_key,
    address_count,
    max_addresses,
    paid_by_time,
    token_issued_at
) VALUES (
    ?1,
    ?2,
    ?3,
    ?4,
    ?5,
    ?6,
    ?7,
    ?8,
    ?9,
    ?10,
    ?11
)
`

type CreateUserParams struct {
	Uuid           string
	Name           string
	IsAdmin        bool
	IsBanned       bool
	Fare           string
	PrivateKey     string
	PublicKey      string
	AddressesCount int64
	MaxAddresses   int64
	PaidByTime     *time.Time
	TokenIssuedAt  *time.Time
}

// Note that generated queries have some parameters marked as nullable, but they really shouldn't.
// In SQLC docs (https://docs.sqlc.dev/en/v1.28.0/howto/named_parameters.html) you can find
// syntax @arg::type, but it doesn't work for me (sqlc reports syntax error).
// Look for real null constraints in comments.
//
//	=========== Users ===========
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, createUser,
		arg.Uuid,
		arg.Name,
		arg.IsAdmin,
		arg.IsBanned,
		arg.Fare,
		arg.PrivateKey,
		arg.PublicKey,
		arg.AddressesCount,
		arg.MaxAddresses,
		arg.PaidByTime,
		arg.TokenIssuedAt,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

const deleteUser = `-- name: DeleteUser :execrows
DELETE FROM users
WHERE uuid = ?1
`

func (q *Queries) DeleteUser(ctx context.Context, uuid string) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteUser, uuid)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT id, uuid, name, is_admin, is_banned, fare, private_key, public_key, address_count, max_addresses, paid_by_time, token_issued_at FROM users
`

func (q *Queries) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.IsAdmin,
			&i.IsBanned,
			&i.Fare,
			&i.PrivateKey,
			&i.PublicKey,
			&i.AddressCount,
			&i.MaxAddresses,
			&i.PaidByTime,
			&i.TokenIssuedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, uuid, name, is_admin, is_banned, fare, private_key, public_key, address_count, max_addresses, paid_by_time, token_issued_at FROM users
WHERE uuid = ?1
`

func (q *Queries) GetUserByName(ctx context.Context, uuid string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, uuid)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.IsAdmin,
		&i.IsBanned,
		&i.Fare,
		&i.PrivateKey,
		&i.PublicKey,
		&i.AddressCount,
		&i.MaxAddresses,
		&i.PaidByTime,
		&i.TokenIssuedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :execrows
UPDATE users SET
    name = ?1,
    is_admin = ?2,
    is_banned = ?3,
    fare = ?4,
    private_key = ?5,
    public_key = ?6,
    address_count = ?7,
    max_addresses = ?8,
    paid_by_time = ?9,
    token_issued_at = ?10
WHERE id = ?11
`

type UpdateUserParams struct {
	Name          string
	IsAdmin       bool
	IsBanned      bool
	Fare          string
	PrivateKey    string
	PublicKey     string
	AddressCount  int64
	MaxAddresses  int64
	PaidByTime    *time.Time
	TokenIssuedAt *time.Time
	ID            int64
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, updateUser,
		arg.Name,
		arg.IsAdmin,
		arg.IsBanned,
		arg.Fare,
		arg.PrivateKey,
		arg.PublicKey,
		arg.AddressCount,
		arg.MaxAddresses,
		arg.PaidByTime,
		arg.TokenIssuedAt,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
