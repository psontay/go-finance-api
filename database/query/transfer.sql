-- name: CreateTransfer :one
insert into transfer(
    from_account_id, to_account_id, amount
) values (
          $1, $2, $3
         ) returning *;

-- name: GetTransfer :one
select * from transfer where id = $1 limit 1;

-- name: ListTransfers :many
select * from transfer order by id
limit  $1
offset $2;
