-- up.sql
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(30) NOT NULL,
    content VARCHAR(255)
);