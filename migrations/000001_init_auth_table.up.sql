-- up.sql
CREATE TABLE IF NOT EXISTS auth (
    id SERIAL PRIMARY KEY,
    login VARCHAR(30) NOT NULL,
    pass_hash TEXT NOT NULL
)