CREATE TABLE IF NOT EXISTS `users`(
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `email` VARCHAR(255) NOT NULL COLLATE NOCASE,
  `phone_number` VARCHAR(255) COLLATE NOCASE,
  `password` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(255) COLLATE NOCASE,
  `last_name` VARCHAR(255) COLLATE NOCASE,
  `lang` VARCHAR(2) DEFAULT 'en' COLLATE NOCASE,
  `role` VARCHAR(255) DEFAULT 'USER' COLLATE NOCASE,
  `super_admin` BOOLEAN DEFAULT 0,
  `status` VARCHAR(255) DEFAULT 'ACTIVE' COLLATE NOCASE,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS `users_email` ON `users`(`email`);
CREATE INDEX IF NOT EXISTS `users_role` ON `users`(`role`);


CREATE TABLE IF NOT EXISTS `otps` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `user_id` INTEGER NOT NULL REFERENCES `users`(`id`) ON DELETE CASCADE,
  `username` VARCHAR(255) NOT NULL COLLATE NOCASE, 
  `code` VARCHAR(255) NOT NULL COLLATE NOCASE,
  `valid_until` TIMESTAMP NOT NULL,
  `max_attempts` INTEGER DEFAULT 3,
  `attempts_count` INTEGER DEFAULT 0,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS `otps_user` ON `otps`(`user_id`);
CREATE INDEX IF NOT EXISTS `otps_username` ON `otps`(`username`);


CREATE TABLE IF NOT EXISTS `assets` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `unique_id` VARCHAR(32) NOT NULL,
  `section` VARCHAR(32) NOT NULL DEFAULT 'default' COLLATE NOCASE,
  `title` VARCHAR(255) COLLATE NOCASE,
  `description` TEXT COLLATE NOCASE,
  `uri` VARCHAR(255) COLLATE NOCASE,
  `storage_name` VARCHAR(255) COLLATE NOCASE,
  `mime_type` VARCHAR(255) COLLATE NOCASE,
  `owner_id` INTEGER,
  `owner_type` VARCHAR(32) COLLATE NOCASE,
  `asset_type` VARCHAR(32) COLLATE NOCASE,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS `assets_owner` ON `assets`(`owner_id`, `owner_type`);
CREATE INDEX IF NOT EXISTS `assets_asset_type` ON `assets`(`asset_type`);
ALTER TABLE `assets`
ADD COLUMN `author_id`
AFTER `id` INTEGER DEFAULT NULL REFERENCES `users`(`id`) ON DELETE
SET NULL;
CREATE INDEX IF NOT EXISTS `assets_author` ON `assets`(`author_id`);
ALTER TABLE `users`
ADD COLUMN `avatar_asset_id` INTEGER DEFAULT NULL REFERENCES `assets`(`id`) ON DELETE
SET NULL;

-- add asset_id to users
ALTER TABLE `users`
ADD COLUMN `asset_id` INTEGER DEFAULT NULL REFERENCES `assets`(`id`) ON DELETE
SET NULL;

-- media library
CREATE TABLE IF NOT EXISTS `asset_libraries` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` VARCHAR(255) NOT NULL COLLATE NOCASE,
  `description` TEXT NOT NULL COLLATE NOCASE,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS `asset_libraries_name` ON `asset_libraries`(`name`);

CREATE TABLE IF NOT EXISTS `asset_folders` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` VARCHAR(255) NOT NULL COLLATE NOCASE,
  `description` TEXT COLLATE NOCASE,
  `path` VARCHAR(512) NOT NULL COLLATE NOCASE,
  `library_id` INTEGER NOT NULL REFERENCES `asset_libraries`(`id`) ON DELETE CASCADE,
  `parent_id` INTEGER REFERENCES `asset_folders`(`id`) ON DELETE SET NULL,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS `asset_folders_name` ON `asset_folders`(`name`, `library_id`);
CREATE UNIQUE INDEX IF NOT EXISTS `asset_folders_path` ON `asset_folders`(`path`, `library_id`);

ALTER TABLE `asset_libraries` 
ADD COLUMN `root_folder_id` INTEGER DEFAULT NULL REFERENCES `asset_folders`(`id`) ON DELETE SET NULL;



-- TAGS
CREATE TABLE IF NOT EXISTS `tags` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `name` VARCHAR(255) NOT NULL COLLATE NOCASE,
  `slug` VARCHAR(255) NOT NULL COLLATE NOCASE,
  `complete_slug` VARCHAR(255) COLLATE NOCASE,
  `owner_type` VARCHAR(32) COLLATE NOCASE,
  `parent_id` INTEGER REFERENCES `tags`(`id`) ON DELETE SET NULL,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX IF NOT EXISTS `tags_owner_type_complete_slug` ON `tags`(`complete_slug`, `owner_type`);
CREATE INDEX IF NOT EXISTS `tags_owner` ON `tags`(`owner_type`);

CREATE TRIGGER IF NOT EXISTS `tags_create_with_parent_complete_slug` 
AFTER INSERT ON `tags`
FOR EACH ROW
WHEN NEW.parent_id IS NOT NULL
BEGIN
 UPDATE tags SET complete_slug = (SELECT complete_slug || '/' || NEW.slug FROM tags WHERE id = NEW.parent_id) WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS `tags_create_without_parent_complete_slug`
AFTER INSERT ON `tags`
FOR EACH ROW
WHEN NEW.parent_id IS NULL
BEGIN
 UPDATE tags SET complete_slug = NEW.slug WHERE id = NEW.id;
END;


-- CREATE TRIGGER IF NOT EXISTS `tags_update_complete_slug`
-- BEFORE UPDATE ON `tags`
-- BEGIN
--   SELECT CASE
--   WHEN NEW.parent_id IS NULL THEN
--     UPDATE tags SET complete_slug = NEW.slug WHERE id = NEW.id;
--   ELSE
--     UPDATE tags SET complete_slug = (SELECT slug || '/' || NEW.slug FROM tags WHERE id = NEW.parent_id) WHERE id = NEW.id;
--   END;
-- END;


CREATE TABLE IF NOT EXISTS `tag_entries` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `tag_id` INTEGER NOT NULL REFERENCES `tags`(`id`) ON DELETE CASCADE,
  `entry_id` INTEGER NOT NULL,
  `entry_type` VARCHAR(32) NOT NULL,
  `metadata` JSON,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS `tag_entries_entry` ON `tag_entries`(`entry_id`, `entry_type`);
CREATE INDEX IF NOT EXISTS `tag_entries_tag` ON `tag_entries`(`tag_id`);

