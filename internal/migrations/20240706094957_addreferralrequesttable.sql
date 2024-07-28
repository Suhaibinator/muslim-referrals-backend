-- +goose Up
-- +goose StatementBegin
CREATE TYPE referral_type AS ENUM ('Internship', 'Full-Time', 'Part-Time', 'Contract');
CREATE TYPE status AS ENUM ('Referral Requested', 'Referred for Job', 'Job Offer Received', 'Job Offer Accepted', 'Job Offer Declined');

CREATE TABLE referral_requests (
    id SERIAL PRIMARY KEY,
    jobTitle VARCHAR(255) NOT NULL,
    jobLinks TEXT[] NOT NULL,
    description TEXT NOT NULL,
    location VARCHAR(255),
    referralType referral_type NOT NULL,
    referrerId INTEGER,
    companyId INTEGER,
    currentStatus status NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updatedAt TIMESTAMP NOT NULL DEFAULT current_timestamp,
    deletedAt TIMESTAMP,
    CONSTRAINT fk_referrer FOREIGN KEY (referrerId) REFERENCES referrers (id) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT fk_company FOREIGN KEY (companyId) REFERENCES companies (id) ON UPDATE CASCADE ON DELETE SET NULL
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE referral_requests;
DROP TYPE referral_type;
DROP TYPE status;
-- +goose StatementEnd
