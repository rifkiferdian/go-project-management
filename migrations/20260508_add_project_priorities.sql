START TRANSACTION;

CREATE TABLE IF NOT EXISTS `project_priorities` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL DEFAULT '#cecece',
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `project_priorities` (`name`, `color`, `is_default`, `created_at`, `updated_at`)
SELECT 'Low', '#008000', 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM `project_priorities` WHERE `name` = 'Low' AND `deleted_at` IS NULL);

INSERT INTO `project_priorities` (`name`, `color`, `is_default`, `created_at`, `updated_at`)
SELECT 'Normal', '#CECECE', 1, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM `project_priorities` WHERE `name` = 'Normal' AND `deleted_at` IS NULL);

INSERT INTO `project_priorities` (`name`, `color`, `is_default`, `created_at`, `updated_at`)
SELECT 'High', '#ff0000', 0, NOW(), NOW()
WHERE NOT EXISTS (SELECT 1 FROM `project_priorities` WHERE `name` = 'High' AND `deleted_at` IS NULL);

ALTER TABLE `projects`
  ADD COLUMN IF NOT EXISTS `priority_id` bigint(20) UNSIGNED DEFAULT NULL AFTER `status_id`;

UPDATE `projects` p
JOIN (
  SELECT MIN(`id`) AS `id`
  FROM `project_priorities`
  WHERE `is_default` = 1
    AND `deleted_at` IS NULL
) pp ON pp.`id` IS NOT NULL
SET p.`priority_id` = pp.`id`
WHERE p.`priority_id` IS NULL
  AND p.`deleted_at` IS NULL;

SET @idx_exists := (
  SELECT COUNT(1)
  FROM information_schema.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'projects'
    AND INDEX_NAME = 'projects_priority_id_foreign'
);
SET @sql := IF(
  @idx_exists = 0,
  'ALTER TABLE `projects` ADD KEY `projects_priority_id_foreign` (`priority_id`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

SET @fk_exists := (
  SELECT COUNT(1)
  FROM information_schema.TABLE_CONSTRAINTS
  WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND TABLE_NAME = 'projects'
    AND CONSTRAINT_NAME = 'projects_priority_id_foreign'
    AND CONSTRAINT_TYPE = 'FOREIGN KEY'
);
SET @sql := IF(
  @fk_exists = 0,
  'ALTER TABLE `projects` ADD CONSTRAINT `projects_priority_id_foreign` FOREIGN KEY (`priority_id`) REFERENCES `project_priorities` (`id`)',
  'SELECT 1'
);
PREPARE stmt FROM @sql;
EXECUTE stmt;
DEALLOCATE PREPARE stmt;

COMMIT;
