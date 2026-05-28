-- name: CreateUser :one
INSERT INTO users (
    username,
    hashed_password,
    full_name,
    email
) VALUES (
             $1, $2, $3, $4
         ) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 LIMIT 1;

-- name: GetUserForUpdate :one
select * from users where username = $1 limit 1 for no key update;

-- name: ListUsers :many
select * from users order by username
limit $1
    offset $2;

-- name: UpdateUser :one
update users
set
    email = $2,
    full_name = $3,
    hashed_password = $4,
    password_changed_at = now()
where username = $1 returning *;

-- name: DeleteUser :exec
delete from users where username = $1;