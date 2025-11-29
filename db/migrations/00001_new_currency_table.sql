-- +goose Up
-- +goose StatementBegin
CREATE users TABLE(
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE INDEX users_email ON users (email);

CREATE tasks TABLE(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    text VARCHAR(255) NOT NULL,
    FOREIGN KEY (userId) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX idx_user_id ON tasks (userId);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
