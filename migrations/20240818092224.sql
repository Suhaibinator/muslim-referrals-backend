-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_companies" table
CREATE TABLE `new_companies` (
  `id` integer NULL PRIMARY KEY AUTOINCREMENT,
  `name` text NOT NULL,
  `domain` text NOT NULL,
  `is_supported` numeric NOT NULL,
  `added_by_user_id` integer NOT NULL,
  `created_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `updated_at` datetime NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  `deleted_at` datetime NULL,
  CONSTRAINT `fk_companies_user` FOREIGN KEY (`added_by_user_id`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Copy rows from old table "companies" to new temporary table "new_companies"
INSERT INTO `new_companies` (`id`, `name`, `domain`, `is_supported`, `created_at`, `updated_at`, `deleted_at`) SELECT `id`, `name`, `domain`, `is_supported`, `created_at`, `updated_at`, `deleted_at` FROM `companies`;
-- Drop "companies" table after copying rows
DROP TABLE `companies`;
-- Rename temporary table "new_companies" to "companies"
ALTER TABLE `new_companies` RENAME TO `companies`;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
