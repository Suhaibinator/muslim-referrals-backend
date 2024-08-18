-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_candidates" table
CREATE TABLE `new_candidates` (
  `candidate_id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `work_experience` integer NOT NULL,
  `resume_url` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL,
  CONSTRAINT `fk_candidates_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Copy rows from old table "candidates" to new temporary table "new_candidates"
INSERT INTO `new_candidates` (`candidate_id`, `user_id`, `work_experience`, `resume_url`, `created_at`) SELECT `candidate_id`, `user_id`, `work_experience`, `resume_url`, `created_at` FROM `candidates`;
-- Drop "candidates" table after copying rows
DROP TABLE `candidates`;
-- Rename temporary table "new_candidates" to "candidates"
ALTER TABLE `new_candidates` RENAME TO `candidates`;
-- Create index "idx_candidates_user_id" to table: "candidates"
CREATE UNIQUE INDEX `idx_candidates_user_id` ON `candidates` (`user_id`, `user_id`);
-- Create "new_referrers" table
CREATE TABLE `new_referrers` (
  `referrer_id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `company_id` integer NOT NULL,
  `corporate_email` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL,
  CONSTRAINT `fk_referrers_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_referrers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "referrers" to new temporary table "new_referrers"
INSERT INTO `new_referrers` (`referrer_id`, `user_id`, `company_id`, `corporate_email`, `created_at`) SELECT `referrer_id`, `user_id`, `company_id`, `corporate_email`, `created_at` FROM `referrers`;
-- Drop "referrers" table after copying rows
DROP TABLE `referrers`;
-- Rename temporary table "new_referrers" to "referrers"
ALTER TABLE `new_referrers` RENAME TO `referrers`;
-- Create index "idx_referrers_user_id" to table: "referrers"
CREATE UNIQUE INDEX `idx_referrers_user_id` ON `referrers` (`user_id`, `user_id`);
-- Create "new_companies" table
CREATE TABLE `new_companies` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `is_supported` numeric NOT NULL,
  `added_by_user_id` integer NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL,
  CONSTRAINT `fk_companies_user` FOREIGN KEY (`added_by_user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "companies" to new temporary table "new_companies"
INSERT INTO `new_companies` (`id`, `name`, `is_supported`, `added_by_user_id`, `created_at`, `updated_at`, `deleted_at`) SELECT `id`, `name`, `is_supported`, `added_by_user_id`, `created_at`, `updated_at`, `deleted_at` FROM `companies`;
-- Drop "companies" table after copying rows
DROP TABLE `companies`;
-- Rename temporary table "new_companies" to "companies"
ALTER TABLE `new_companies` RENAME TO `companies`;
-- Create "company_domain_associations" table
CREATE TABLE `company_domain_associations` (
  `company_id` integer NULL,
  `domain` text NULL,
  PRIMARY KEY (`company_id`, `domain`),
  CONSTRAINT `fk_companies_domains` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
