-- +goose Up
-- +goose StatementBegin
CREATE TABLE referrers (
    id SERIAL PRIMARY KEY,
    corporateEmail VARCHAR(255) NOT NULL,
    companyId INTEGER,
    userId INTEGER,
    createdAt TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updatedAt TIMESTAMP NOT NULL DEFAULT current_timestamp,
    deletedAt TIMESTAMP,
    CONSTRAINT fk_company FOREIGN KEY (companyId) REFERENCES companies (id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_user FOREIGN KEY (userId) REFERENCES users (id) ON UPDATE CASCADE ON DELETE SET NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE referrers;
-- +goose StatementEnd
