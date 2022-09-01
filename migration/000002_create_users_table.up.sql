CREATE TABLE IF NOT EXISTS users(
    id SERIAL NOT NULL PRIMARY KEY,
    email VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);