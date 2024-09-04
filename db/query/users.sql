-- name: GetUserRole :one
select users.username, users.password, users.email, group_concat(role_details.detail) as role_details from users
inner join user_roles on users.id = user_roles.user_id
inner join role_details on user_roles.role_id = role_details.id
where username = ?
group by users.id
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


