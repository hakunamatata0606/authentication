-- name: GetUserRole :one
select users.username, users.password, users.email, role_details.detail as role_detail from users
inner join user_roles on users.id = user_roles.user_id
inner join role_details on user_roles.role_id = role_details.id
where username = ?
limit 1;

-- name: GetUserNoUpdate :one
select * from users
where username = ? limit 1;

-- name: GetUserForUpdate :one
select * from users
where username = ? limit 1
for update;

-- name: ListUser :many
select * from users;

-- name: CreateUser :execresult
insert into users (username, password, email)
values (?, ?, ?);

-- name: DeleteUser :execresult
delete from users
where username = ?;


