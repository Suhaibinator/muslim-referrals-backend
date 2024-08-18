-- Create "users" table
CREATE TABLE `users` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `first_name` text NOT NULL,
  `last_name` text NOT NULL,
  `email` text NOT NULL,
  `phone_number` text NULL,
  `phone_ext` text NULL,
  `linked_in` text NULL,
  `github` text NULL,
  `website` text NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL
);
-- Create index "idx_users_email" to table: "users"
CREATE UNIQUE INDEX `idx_users_email` ON `users` (`email`);
-- Create "candidates" table
CREATE TABLE `candidates` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `work_experience` integer NOT NULL,
  `resume_url` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT `fk_candidates_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_candidates_user_id" to table: "candidates"
CREATE UNIQUE INDEX `idx_candidates_user_id` ON `candidates` (`user_id`);
-- Create "companies" table
CREATE TABLE `companies` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `domain` text NOT NULL,
  `is_supported` numeric NOT NULL,
  `location` text NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL
);
-- Create "referrers" table
CREATE TABLE `referrers` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `user_id` integer NOT NULL,
  `company_id` integer NOT NULL,
  `corporate_email` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  CONSTRAINT `fk_referrers_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_referrers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_referrers_user_id" to table: "referrers"
CREATE UNIQUE INDEX `idx_referrers_user_id` ON `referrers` (`user_id`);
-- Create "referral_requests" table
CREATE TABLE `referral_requests` (
  `candidate_id` integer NULL,
  `company_id` integer NULL,
  `primary_job_title_seeking` text NOT NULL,
  `summary` text NOT NULL,
  `referral_type` text NOT NULL,
  `referrer_id` integer NULL,
  `status` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL,
  PRIMARY KEY (`candidate_id`, `company_id`),
  CONSTRAINT `fk_referral_requests_company` FOREIGN KEY (`company_id`) REFERENCES `companies` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_referral_requests_candidate` FOREIGN KEY (`candidate_id`) REFERENCES `candidates` (`id`) ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT `fk_referral_requests_referrer` FOREIGN KEY (`referrer_id`) REFERENCES `referrers` (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create "referral_request_job_links_associations" table
CREATE TABLE `referral_request_job_links_associations` (
  `referral_request_id` integer NULL,
  `job_link` text NULL,
  PRIMARY KEY (`referral_request_id`, `job_link`),
  CONSTRAINT `fk_referral_requests_job_links` FOREIGN KEY (`referral_request_id`) REFERENCES `referral_requests` (`candidate_id`) ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create "referral_request_location_associations" table
CREATE TABLE `referral_request_location_associations` (
  `referral_request_id` integer NULL,
  `location` text NULL,
  PRIMARY KEY (`referral_request_id`, `location`),
  CONSTRAINT `fk_referral_requests_locations` FOREIGN KEY (`referral_request_id`) REFERENCES `referral_requests` (`candidate_id`) ON UPDATE CASCADE ON DELETE CASCADE
);
