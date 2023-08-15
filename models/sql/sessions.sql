CREATE TABLE sessions (
    id SERIAL PRIMARY KEY, 
    user_id INT UNIQUE,
    token_hash TEXT UNIQUE NOT NULL
);

-- user_id will be used to create a relationship between a session and a specific user