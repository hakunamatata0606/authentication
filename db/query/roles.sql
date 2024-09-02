-- name: AddUserRole :execresult
insert into user_roles (user_id, role_id)
values (? , ?);

-- name: GetRoleIdByDetail :one
select * from role_details
where detail = ?
limit 1;

-- name: DeleteRole :execresult
delete from user_roles
where user_id = ?;
