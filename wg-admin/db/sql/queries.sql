-- Note that generated queries have some parameters marked as nullable, but they really shouldn't.
-- In SQLC docs (https://docs.sqlc.dev/en/v1.28.0/howto/named_parameters.html) you can find 
-- syntax @arg::type, but it doesn't work for me (sqlc reports syntax error).
-- Look for real null constraints in comments.

--  =========== Users ===========

-- name: CreateUser :execlastid
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
    token_issued_at,
    last_seen_at
) VALUES (
    @uuid,
    @name,
    @is_admin,
    @is_banned,
    @fare,
    @private_key,
    @public_key,
    @addresses_count,
    @max_addresses,
    @paid_by_time,
    @token_issued_at,
    @last_seen_at
);

-- name: GetUserByName :one
SELECT * FROM users
WHERE name = @name;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: UpdateUser :execrows
UPDATE users SET
    name = @name,
    is_admin = @is_admin,
    is_banned = @is_banned,
    fare = @fare,
    private_key = @private_key,
    public_key = @public_key,
    address_count = @address_count,
    max_addresses = @max_addresses,
    paid_by_time = @paid_by_time,
    token_issued_at = @token_issued_at,
    last_seen_at = @last_seen_at
WHERE id = @id;

-- name: DeleteUser :execrows
DELETE FROM users
WHERE name = @name;

--  =========== Addresses ===========

-- -- name: GetAddressCount :one
-- SELECT host_id from addresses
-- WHERE is_v6 = @is_v6
-- ORDER BY host_id DESC
-- LIMIT 1;
-- 
-- -- name: FindLowestFreeAddress :one
-- SELECT * from addresses
-- WHERE is_v6 = @is_v6 AND user_id IS NOT NULL
-- ORDER BY host_id DESC
-- LIMIT 1;
-- 
-- -- name: CreateAddress :one
-- INSERT INTO addresses (
--     host_id,
--     name,
--     is_v6,
--     user_id,
--     desynced_at
-- )
-- SELECT 
--     host_id + 1,
--     @new_name,
--     @is_v6,
--     @user_id,
--     @desynced_at
-- FROM addresses
-- WHERE
--     is_v6 = @is_v6
-- ORDER BY host_id DESC
-- LIMIT 1
-- RETURNING *;
-- 
-- -- name: UpdateAddress :one
-- UPDATE addresses SET
--     name = @name,
--     is_v6 = @is_v6,
--     user_id = @user_id,
--     desynced_at = @desynced_at
-- WHERE host_id = @host_id
-- RETURNING *;
-- 
-- -- name: GetAddress :one
-- SELECT addresses.*, users.* FROM addresses
-- JOIN users ON users.id = addresses.user_id
-- WHERE host_id = @host_id;
-- 
-- -- name: GetAddressesOwnedBy :many
-- SELECT addresses.*, users.* FROM addresses
-- JOIN users ON users.id = addresses.user_id
-- WHERE user_id = @user_id;
-- 
-- 
-- -- Cannot use sqlc.embed() here, since it doesn't make user nullable.
-- -- See https://github.com/sqlc-dev/sqlc/issues/3240
-- -- name: GetAddressesDesynced :many
-- SELECT addresses.*, users.* FROM addresses
-- LEFT OUTER JOIN users ON users.id = addresses.user_id
-- WHERE
--     desynced_at > @not_before AND
--     is_v6 = @is_v6 AND
--     host_id < @max_host_id;
