START TRANSACTION;

CREATE TABLE IF NOT EXISTS `ticket_template_epics` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `set_id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` longtext DEFAULT NULL,
  `start_offset_days` int(11) DEFAULT NULL,
  `due_offset_days` int(11) DEFAULT NULL,
  `sort_order` int(11) NOT NULL DEFAULT 1,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ticket_template_epics_set_name_unique` (`set_id`, `name`),
  KEY `ticket_template_epics_set_id_index` (`set_id`),
  KEY `ticket_template_epics_sort_order_index` (`set_id`, `sort_order`),
  CONSTRAINT `ticket_template_epics_set_id_foreign` FOREIGN KEY (`set_id`) REFERENCES `ticket_template_sets` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

ALTER TABLE `ticket_template_items`
  ADD COLUMN IF NOT EXISTS `template_epic_id` bigint(20) UNSIGNED DEFAULT NULL AFTER `description`;

SET @idx_exists := (
  SELECT COUNT(1)
  FROM INFORMATION_SCHEMA.STATISTICS
  WHERE TABLE_SCHEMA = DATABASE()
    AND TABLE_NAME = 'ticket_template_items'
    AND INDEX_NAME = 'ticket_template_items_template_epic_id_index'
);
SET @idx_sql := IF(
  @idx_exists = 0,
  'ALTER TABLE `ticket_template_items` ADD KEY `ticket_template_items_template_epic_id_index` (`template_epic_id`)',
  'SELECT 1'
);
PREPARE idx_stmt FROM @idx_sql;
EXECUTE idx_stmt;
DEALLOCATE PREPARE idx_stmt;

SET @fk_exists := (
  SELECT COUNT(1)
  FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
  WHERE CONSTRAINT_SCHEMA = DATABASE()
    AND CONSTRAINT_NAME = 'ticket_template_items_template_epic_id_foreign'
);
SET @fk_sql := IF(
  @fk_exists = 0,
  'ALTER TABLE `ticket_template_items` ADD CONSTRAINT `ticket_template_items_template_epic_id_foreign` FOREIGN KEY (`template_epic_id`) REFERENCES `ticket_template_epics` (`id`) ON DELETE SET NULL',
  'SELECT 1'
);
PREPARE fk_stmt FROM @fk_sql;
EXECUTE fk_stmt;
DEALLOCATE PREPARE fk_stmt;

COMMIT;
