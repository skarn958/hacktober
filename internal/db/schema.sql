CREATE TABLE IF NOT EXISTS users (
    fingerprint TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS solves (
    fingerprint  TEXT NOT NULL REFERENCES users(fingerprint) ON DELETE CASCADE,
    challenge_id TEXT NOT NULL,
    solved_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (fingerprint, challenge_id)
);