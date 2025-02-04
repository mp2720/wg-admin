-- one user <-> many addresses

CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid            TEXT NOT NULL UNIQUE,
    name            TEXT NOT NULL UNIQUE,
    is_admin        BOOLEAN NOT NULL,
    is_banned       BOOLEAN NOT NULL,
    fare            TEXT NOT NULL,
    private_key     TEXT NOT NULL,
    public_key      TEXT NOT NULL UNIQUE,
    address_count   INTEGER NOT NULL,
    max_addresses   INTEGER NOT NULL,

    paid_by_time    TIMESTAMP,
    token_issued_at TIMESTAMP,
    last_seen_at    TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_uuid ON users(uuid);
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_name ON users(name);

-- when updating config, make sure no other processes are working with db and the net masks are
-- big enough for all assigned addresses, so the data won't be corrupted
CREATE TABLE IF NOT EXISTS wg_net_config (
    id      INTEGER PRIMARY KEY CHECK (id = 0), -- allow <= 1 rows
    v6_net  TEXT NOT NULL,
    v4_net  TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS addresses (
    id              INTEGER PRIMARY KEY,
    is_reserved     BOOLEAN NOT NULL,
    is_v6           BOOLEAN NOT NULL,
    host_id         INTEGER NOT NULL, -- unique for all rows with the same version
    owner_user_id   INTEGER,
    description     TEXT,
    desynced_at     TIMESTAMP NOT NULL,

    FOREIGN KEY(owner_user_id) REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_addresses_owne_user_id ON addresses(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_addresses_host_id ON addresses(host_id);

CREATE TRIGGER IF NOT EXISTS trigger_abort_delete_addresses
BEFORE DELETE ON addresses
BEGIN
    SELECT RAISE(ABORT, 'address deletion is prohibited');
END;

