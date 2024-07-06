-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    fullName VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phoneNumber VARCHAR(255),
    resumeUrl VARCHAR(255),
    linkedIn VARCHAR(255),
    github VARCHAR(255),
    website VARCHAR(255),
    createdAt TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updatedAt TIMESTAMP NOT NULL DEFAULT current_timestamp,
    deletedAt TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
