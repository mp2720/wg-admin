-- User owns many addresses and statistics rows.
-- Server owns many statistics rows.

--  =========== User ===========
CREATE TABLE IF NOT EXISTS users (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
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

--  =========== Address ===========
CREATE TABLE IF NOT EXISTS addresses (
    id          INTEGER PRIMARY KEY, -- AUTOINCREMENT is not required, since deletion is prohibited
    host_id     INTEGER NOT NULL,
    name        TEXT NOT NULL,
    is_v6       BOOLEAN NOT NULL,
    user_id     INTEGER,
    desynced_at TIMESTAMP NOT NULL,

    FOREIGN KEY(user_id) REFERENCES users(id)
        ON DELETE RESTRICT
        ON UPDATE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_addresses_user_id ON addresses(user_id);

CREATE TRIGGER IF NOT EXISTS trigger_abort_delete_address
BEFORE DELETE ON addresses
BEGIN
    SELECT RAISE(ABORT, 'address deletion is prohibited');
END;

