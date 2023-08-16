CREATE TABLE sessions (
    id SERIAL PRIMARY KEY,
    user_id INT UNIQUE REFERENCES users (id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL
);
-- user_id will be used to create a relationship between a session and a specific user
SELECT users.id,
    users.email,
    users.password_hash
FROM sessions
    JOIN users ON users.id = sessions.user_id
WHERE sessions.token_hash = $1;

