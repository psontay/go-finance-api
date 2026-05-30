-- name: CreateTransfer :one
insert into transfer(
    from_account_id, to_account_id, amount
) values (
          $1, $2, $3
         ) returning *;

-- name: GetTransfer :one
select * from transfer where id = $1 limit 1;

-- name: ListTransfers :many
SELECT * FROM transfer
WHERE from_account_id = $1 OR to_account_id = $2
ORDER BY id
LIMIT $3 OFFSET $4;
