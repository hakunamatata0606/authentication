// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package db

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :execresult
insert into users (username, password, email)
values (?, ?, ?)
`

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createUser, arg.Username, arg.Password, arg.Email)
}

const deleteUser = `-- name: DeleteUser :execresult
delete from users
where username = ?
`

func (q *Queries) DeleteUser(ctx context.Context, username string) (sql.Result, error) {
	return q.db.ExecContext(ctx, deleteUser, username)
}

const getUserForUpdate = `-- name: GetUserForUpdate :one
select id, username, password, email from users
where username = ? limit 1
for update
`

func (q *Queries) GetUserForUpdate(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserForUpdate, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const getUserNoUpdate = `-- name: GetUserNoUpdate :one
select id, username, password, email from users
where username = ? limit 1
`

func (q *Queries) GetUserNoUpdate(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserNoUpdate, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
	)
	return i, err
}

const getUserRole = `-- name: GetUserRole :one
select users.username, users.password, users.email, role_details.detail as role_detail from users
inner join user_roles on users.id = user_roles.user_id
inner join role_details on user_roles.role_id = role_details.id
where username = ?
limit 1
`

type GetUserRoleRow struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	Email      string `json:"email"`
	RoleDetail string `json:"role_detail"`
}

func (q *Queries) GetUserRole(ctx context.Context, username string) (GetUserRoleRow, error) {
	row := q.db.QueryRowContext(ctx, getUserRole, username)
	var i GetUserRoleRow
	err := row.Scan(
		&i.Username,
		&i.Password,
		&i.Email,
		&i.RoleDetail,
	)
	return i, err
}

const listUser = `-- name: ListUser :many
select id, username, password, email from users
`

func (q *Queries) ListUser(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Password,
			&i.Email,
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
