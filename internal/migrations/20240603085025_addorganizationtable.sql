-- +goose Up
-- +goose StatementBegin
CREATE TABLE ORGANIZATION (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    domain VARCHAR(255) NOT NULL,
    is_supported BOOLEAN NOT NULL,
    deleted_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE ORGANIZATION;
-- +goose StatementEnd
