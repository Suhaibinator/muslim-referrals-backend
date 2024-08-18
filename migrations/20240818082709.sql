-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Rename a column from "id" to "candidate_id"
ALTER TABLE `candidates` RENAME COLUMN `id` TO `candidate_id`;
-- Create "new_companies" table
CREATE TABLE `new_companies` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `domain` text NOT NULL,
  `is_supported` numeric NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL
);
-- Copy rows from old table "companies" to new temporary table "new_companies"
INSERT INTO `new_companies` (`id`, `name`, `domain`, `is_supported`, `created_at`, `updated_at`, `deleted_at`) SELECT `id`, `name`, `domain`, `is_supported`, `created_at`, `updated_at`, `deleted_at` FROM `companies`;
-- Drop "companies" table after copying rows
DROP TABLE `companies`;
-- Rename temporary table "new_companies" to "companies"
ALTER TABLE `new_companies` RENAME TO `companies`;
-- Rename a column from "id" to "referrer_id"
ALTER TABLE `referrers` RENAME COLUMN `id` TO `referrer_id`;
-- Create "new_referral_requests" table
CREATE TABLE `new_referral_requests` (
  `referral_request_id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `candidate_id` integer NOT NULL,
  `company_id` integer NOT NULL,
  `primary_job_title_seeking` text NOT NULL,
  `summary` text NOT NULL,
  `referral_type` text NOT NULL,
  `referrer_id` integer NULL,
  `status` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL,
  CONSTRAINT `fk_referral_requests_referrer` FOREIGN KEY (`referrer_id`) REFERENCES `referrers` (`referrer_id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_referral_requests_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_referral_requests_candidate` FOREIGN KEY (`candidate_id`) REFERENCES `candidates` (`candidate_id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "referral_requests" to new temporary table "new_referral_requests"
INSERT INTO `new_referral_requests` (`candidate_id`, `company_id`, `primary_job_title_seeking`, `summary`, `referral_type`, `referrer_id`, `status`, `created_at`, `updated_at`, `deleted_at`) SELECT `candidate_id`, `company_id`, `primary_job_title_seeking`, `summary`, `referral_type`, `referrer_id`, `status`, `created_at`, `updated_at`, `deleted_at` FROM `referral_requests`;
-- Drop "referral_requests" table after copying rows
DROP TABLE `referral_requests`;
-- Rename temporary table "new_referral_requests" to "referral_requests"
ALTER TABLE `new_referral_requests` RENAME TO `referral_requests`;
-- Create "new_referral_request_job_links_associations" table
CREATE TABLE `new_referral_request_job_links_associations` (
  `referral_request_id` integer NULL,
  `job_link` text NULL,
  PRIMARY KEY (`referral_request_id`, `job_link`),
  CONSTRAINT `fk_referral_requests_job_links` FOREIGN KEY (`referral_request_id`) REFERENCES `referral_requests` (`referral_request_id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "referral_request_job_links_associations" to new temporary table "new_referral_request_job_links_associations"
INSERT INTO `new_referral_request_job_links_associations` (`referral_request_id`, `job_link`) SELECT `referral_request_id`, `job_link` FROM `referral_request_job_links_associations`;
-- Drop "referral_request_job_links_associations" table after copying rows
DROP TABLE `referral_request_job_links_associations`;
-- Rename temporary table "new_referral_request_job_links_associations" to "referral_request_job_links_associations"
ALTER TABLE `new_referral_request_job_links_associations` RENAME TO `referral_request_job_links_associations`;
-- Create "new_referral_request_location_associations" table
CREATE TABLE `new_referral_request_location_associations` (
  `referral_request_id` integer NULL,
  `location` text NULL,
  PRIMARY KEY (`referral_request_id`, `location`),
  CONSTRAINT `fk_referral_requests_locations` FOREIGN KEY (`referral_request_id`) REFERENCES `referral_requests` (`referral_request_id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "referral_request_location_associations" to new temporary table "new_referral_request_location_associations"
INSERT INTO `new_referral_request_location_associations` (`referral_request_id`, `location`) SELECT `referral_request_id`, `location` FROM `referral_request_location_associations`;
-- Drop "referral_request_location_associations" table after copying rows
DROP TABLE `referral_request_location_associations`;
-- Rename temporary table "new_referral_request_location_associations" to "referral_request_location_associations"
ALTER TABLE `new_referral_request_location_associations` RENAME TO `referral_request_location_associations`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
