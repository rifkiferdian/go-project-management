START TRANSACTION;

CREATE TABLE IF NOT EXISTS `project_divisions` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_id` bigint(20) UNSIGNED NOT NULL,
  `division_id` bigint(20) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `project_divisions_project_division_unique` (`project_id`,`division_id`),
  KEY `project_divisions_project_id_foreign` (`project_id`),
  KEY `project_divisions_division_id_foreign` (`division_id`),
  CONSTRAINT `project_divisions_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`) ON DELETE CASCADE,
  CONSTRAINT `project_divisions_division_id_foreign` FOREIGN KEY (`division_id`) REFERENCES `divisions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Backfill data existing: requester division mengikuti divisi owner project.
INSERT INTO `project_divisions` (`project_id`, `division_id`, `created_at`, `updated_at`)
SELECT
  p.id AS project_id,
  ud.division_id,
  NOW(),
  NOW()
FROM `projects` p
JOIN `user_divisions` ud ON ud.user_id = p.owner_id
LEFT JOIN `project_divisions` pd
  ON pd.project_id = p.id
  AND pd.division_id = ud.division_id
WHERE p.deleted_at IS NULL
  AND pd.id IS NULL;

COMMIT;
