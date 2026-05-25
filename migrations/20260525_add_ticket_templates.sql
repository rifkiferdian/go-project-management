START TRANSACTION;

CREATE TABLE IF NOT EXISTS `ticket_template_sets` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(150) NOT NULL,
  `purpose` varchar(50) NOT NULL COMMENT 'new_project|new_feature|bugfix|custom',
  `description` varchar(255) DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ticket_template_sets_name_purpose_unique` (`name`, `purpose`),
  KEY `ticket_template_sets_purpose_index` (`purpose`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

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

CREATE TABLE IF NOT EXISTS `ticket_template_items` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `set_id` bigint(20) UNSIGNED NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` longtext DEFAULT NULL,
  `template_epic_id` bigint(20) UNSIGNED DEFAULT NULL,
  `default_type_id` bigint(20) UNSIGNED DEFAULT NULL,
  `default_priority_id` bigint(20) UNSIGNED DEFAULT NULL,
  `default_status_id` bigint(20) UNSIGNED DEFAULT NULL,
  `default_owner_id` bigint(20) UNSIGNED DEFAULT NULL,
  `default_responsible_id` bigint(20) UNSIGNED DEFAULT NULL,
  `estimation` double(8,2) DEFAULT NULL,
  `start_offset_days` int(11) DEFAULT NULL,
  `due_offset_days` int(11) DEFAULT NULL,
  `sort_order` int(11) NOT NULL DEFAULT 1,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `ticket_template_items_set_id_index` (`set_id`),
  KEY `ticket_template_items_template_epic_id_index` (`template_epic_id`),
  KEY `ticket_template_items_sort_order_index` (`set_id`, `sort_order`),
  KEY `ticket_template_items_default_type_id_index` (`default_type_id`),
  KEY `ticket_template_items_default_priority_id_index` (`default_priority_id`),
  KEY `ticket_template_items_default_status_id_index` (`default_status_id`),
  KEY `ticket_template_items_default_owner_id_index` (`default_owner_id`),
  KEY `ticket_template_items_default_responsible_id_index` (`default_responsible_id`),
  CONSTRAINT `ticket_template_items_set_id_foreign` FOREIGN KEY (`set_id`) REFERENCES `ticket_template_sets` (`id`) ON DELETE CASCADE,
  CONSTRAINT `ticket_template_items_template_epic_id_foreign` FOREIGN KEY (`template_epic_id`) REFERENCES `ticket_template_epics` (`id`) ON DELETE SET NULL,
  CONSTRAINT `ticket_template_items_default_type_id_foreign` FOREIGN KEY (`default_type_id`) REFERENCES `ticket_types` (`id`) ON DELETE SET NULL,
  CONSTRAINT `ticket_template_items_default_priority_id_foreign` FOREIGN KEY (`default_priority_id`) REFERENCES `ticket_priorities` (`id`) ON DELETE SET NULL,
  CONSTRAINT `ticket_template_items_default_status_id_foreign` FOREIGN KEY (`default_status_id`) REFERENCES `ticket_statuses` (`id`) ON DELETE SET NULL,
  CONSTRAINT `ticket_template_items_default_owner_id_foreign` FOREIGN KEY (`default_owner_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `ticket_template_items_default_responsible_id_foreign` FOREIGN KEY (`default_responsible_id`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

COMMIT;
