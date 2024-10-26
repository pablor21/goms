CREATE TABLE IF NOT EXISTS `users`(
  `id` INTEGER PRIMARY KEY AUTOINCREMENT,
  `email` VARCHAR(255) NOT NULL,
  `phone_number` VARCHAR(255),
  `password` VARCHAR(255) NOT NULL,
  `first_name` VARCHAR(255),
  `last_name` VARCHAR(255),
  `lang` VARCHAR(2) DEFAULT 'en',
  `role` VARCHAR(255) DEFAULT 'USER',
  `super_admin` BOOLEAN DEFAULT 0,
  `status` VARCHAR(255) DEFAULT 'ACTIVE',
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE UNIQUE INDEX IF NOT EXISTS `users_email` ON `users`(`email`);
CREATE INDEX IF NOT EXISTS `users_role` ON `users`(`role`);