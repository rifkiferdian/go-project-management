START TRANSACTION;

CREATE TABLE IF NOT EXISTS `approval_flows` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `flow_code` varchar(100) NOT NULL,
  `flow_name` varchar(255) NOT NULL,
  `entity_type` enum('project_request') NOT NULL DEFAULT 'project_request',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `approval_flows_flow_code_unique` (`flow_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `approval_flow_steps` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `approval_flow_id` bigint(20) UNSIGNED NOT NULL,
  `step_order` int(10) UNSIGNED NOT NULL,
  `step_name` varchar(255) NOT NULL,
  `approval_rule` enum('any','all') NOT NULL DEFAULT 'any',
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `approval_flow_steps_flow_order_unique` (`approval_flow_id`,`step_order`),
  KEY `approval_flow_steps_approval_flow_id_foreign` (`approval_flow_id`),
  CONSTRAINT `approval_flow_steps_approval_flow_id_foreign` FOREIGN KEY (`approval_flow_id`) REFERENCES `approval_flows` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `approval_flow_step_approvers` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `approval_flow_step_id` bigint(20) UNSIGNED NOT NULL,
  `approver_type` enum('user','role','division') NOT NULL,
  `approver_user_id` bigint(20) UNSIGNED DEFAULT NULL,
  `approver_role_id` bigint(20) UNSIGNED DEFAULT NULL,
  `approver_division_id` bigint(20) UNSIGNED DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT 1,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `approval_flow_step_approvers_flow_step_id_foreign` (`approval_flow_step_id`),
  KEY `approval_flow_step_approvers_user_id_foreign` (`approver_user_id`),
  KEY `approval_flow_step_approvers_role_id_foreign` (`approver_role_id`),
  KEY `approval_flow_step_approvers_division_id_foreign` (`approver_division_id`),
  CONSTRAINT `approval_flow_step_approvers_flow_step_id_foreign` FOREIGN KEY (`approval_flow_step_id`) REFERENCES `approval_flow_steps` (`id`) ON DELETE CASCADE,
  CONSTRAINT `approval_flow_step_approvers_user_id_foreign` FOREIGN KEY (`approver_user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `approval_flow_step_approvers_role_id_foreign` FOREIGN KEY (`approver_role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `approval_flow_step_approvers_division_id_foreign` FOREIGN KEY (`approver_division_id`) REFERENCES `divisions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_requests` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `request_no` varchar(50) DEFAULT NULL,
  `project_name` varchar(255) NOT NULL,
  `project_description` text DEFAULT NULL,
  `business_goal` text DEFAULT NULL,
  `request_division_id` bigint(20) UNSIGNED NOT NULL,
  `requested_ticket_prefix` varchar(3) NOT NULL,
  `requester_name` varchar(255) NOT NULL,
  `requester_email` varchar(255) NOT NULL,
  `requester_phone` varchar(50) DEFAULT NULL,
  `requester_employee_id` varchar(100) DEFAULT NULL,
  `approval_flow_id` bigint(20) UNSIGNED NOT NULL,
  `current_step_order` int(10) UNSIGNED DEFAULT 1,
  `status` enum('pending','approved','rejected','synced_to_project','cancelled') NOT NULL DEFAULT 'pending',
  `final_decided_by` bigint(20) UNSIGNED DEFAULT NULL,
  `final_decided_at` timestamp NULL DEFAULT NULL,
  `rejection_reason` text DEFAULT NULL,
  `project_id` bigint(20) UNSIGNED DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `project_requests_request_no_unique` (`request_no`),
  KEY `project_requests_request_division_id_foreign` (`request_division_id`),
  KEY `project_requests_approval_flow_id_foreign` (`approval_flow_id`),
  KEY `project_requests_project_id_foreign` (`project_id`),
  KEY `project_requests_status_idx` (`status`),
  CONSTRAINT `project_requests_request_division_id_foreign` FOREIGN KEY (`request_division_id`) REFERENCES `divisions` (`id`),
  CONSTRAINT `project_requests_approval_flow_id_foreign` FOREIGN KEY (`approval_flow_id`) REFERENCES `approval_flows` (`id`),
  CONSTRAINT `project_requests_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_request_step_states` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_request_id` bigint(20) UNSIGNED NOT NULL,
  `approval_flow_step_id` bigint(20) UNSIGNED NOT NULL,
  `step_order` int(10) UNSIGNED NOT NULL,
  `step_name` varchar(255) NOT NULL,
  `approval_rule` enum('any','all') NOT NULL,
  `status` enum('pending','approved','rejected','skipped') NOT NULL DEFAULT 'pending',
  `decided_by` bigint(20) UNSIGNED DEFAULT NULL,
  `decided_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `project_request_step_states_request_step_unique` (`project_request_id`,`approval_flow_step_id`),
  KEY `project_request_step_states_request_id_foreign` (`project_request_id`),
  KEY `project_request_step_states_flow_step_id_foreign` (`approval_flow_step_id`),
  KEY `project_request_step_states_decided_by_foreign` (`decided_by`),
  CONSTRAINT `project_request_step_states_request_id_foreign` FOREIGN KEY (`project_request_id`) REFERENCES `project_requests` (`id`) ON DELETE CASCADE,
  CONSTRAINT `project_request_step_states_flow_step_id_foreign` FOREIGN KEY (`approval_flow_step_id`) REFERENCES `approval_flow_steps` (`id`),
  CONSTRAINT `project_request_step_states_decided_by_foreign` FOREIGN KEY (`decided_by`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `project_request_approvals` (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `project_request_id` bigint(20) UNSIGNED NOT NULL,
  `approval_flow_step_id` bigint(20) UNSIGNED NOT NULL,
  `approver_user_id` bigint(20) UNSIGNED NOT NULL,
  `decision` enum('approved','rejected') NOT NULL,
  `note` text DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `project_request_approvals_unique` (`project_request_id`,`approval_flow_step_id`,`approver_user_id`),
  KEY `project_request_approvals_request_id_foreign` (`project_request_id`),
  KEY `project_request_approvals_flow_step_id_foreign` (`approval_flow_step_id`),
  KEY `project_request_approvals_approver_user_id_foreign` (`approver_user_id`),
  CONSTRAINT `project_request_approvals_request_id_foreign` FOREIGN KEY (`project_request_id`) REFERENCES `project_requests` (`id`) ON DELETE CASCADE,
  CONSTRAINT `project_request_approvals_flow_step_id_foreign` FOREIGN KEY (`approval_flow_step_id`) REFERENCES `approval_flow_steps` (`id`),
  CONSTRAINT `project_request_approvals_approver_user_id_foreign` FOREIGN KEY (`approver_user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'List project requests', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'List project requests' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Approve project request', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Approve project request' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Reject project request', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Reject project request' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'List approval flows', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'List approval flows' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Create approval flow', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Create approval flow' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Update approval flow', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Update approval flow' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Delete approval flow', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Delete approval flow' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'List approval flow steps', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'List approval flow steps' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Create approval flow step', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Create approval flow step' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Update approval flow step', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Update approval flow step' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Delete approval flow step', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Delete approval flow step' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'List approval flow step approvers', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'List approval flow step approvers' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Create approval flow step approver', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Create approval flow step approver' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Update approval flow step approver', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Update approval flow step approver' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Delete approval flow step approver', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Delete approval flow step approver' AND `guard_name` = 'web'
);

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'List projects' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'List project requests'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Update project' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Approve project request'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Update project' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Reject project request'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'List roles' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'List approval flows'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Create role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Create approval flow'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Update role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Update approval flow'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Delete role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Delete approval flow'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'List roles' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'List approval flow steps'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Create role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Create approval flow step'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Update role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Update approval flow step'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Delete role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Delete approval flow step'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'List roles' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'List approval flow step approvers'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Create role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Create approval flow step approver'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Update role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Update approval flow step approver'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT p_new.id, rhp.role_id
FROM `permissions` p_new
JOIN `permissions` p_old
  ON p_old.name = 'Delete role' AND p_old.guard_name = 'web'
JOIN `role_has_permissions` rhp
  ON rhp.permission_id = p_old.id
LEFT JOIN `role_has_permissions` x
  ON x.permission_id = p_new.id AND x.role_id = rhp.role_id
WHERE p_new.name = 'Delete approval flow step approver'
  AND p_new.guard_name = 'web'
  AND x.permission_id IS NULL;

COMMIT;

-- =============================================================
-- Contoh setup flow approval dinamis
-- =============================================================

-- A) Buat 1 flow aktif untuk project request
/*
INSERT INTO `approval_flows`
  (`flow_code`, `flow_name`, `entity_type`, `is_active`, `created_at`, `updated_at`)
VALUES
  ('PRJ-DEFAULT', 'Project Request Default Flow', 'project_request', 1, NOW(), NOW());
*/

-- B) Buat step (bisa bebas jumlah dan urutan)
/*
INSERT INTO `approval_flow_steps`
  (`approval_flow_id`, `step_order`, `step_name`, `approval_rule`, `is_active`, `created_at`, `updated_at`)
VALUES
  (1, 1, 'Review Divisi Requester', 'any', 1, NOW(), NOW()),
  (1, 2, 'Review IT Governance', 'all', 1, NOW(), NOW()),
  (1, 3, 'Review Finance', 'any', 1, NOW(), NOW());
*/

-- C) Set approver per step (fleksibel: user/role/division)
/*
-- step 1: role manager divisi
INSERT INTO `approval_flow_step_approvers`
  (`approval_flow_step_id`, `approver_type`, `approver_role_id`, `is_active`, `created_at`, `updated_at`)
VALUES
  (1, 'role', 3, 1, NOW(), NOW());

-- step 2: user spesifik IT lead + IT manager (karena rule = all, dua-duanya harus approve)
INSERT INTO `approval_flow_step_approvers`
  (`approval_flow_step_id`, `approver_type`, `approver_user_id`, `is_active`, `created_at`, `updated_at`)
VALUES
  (2, 'user', 7, 1, NOW(), NOW()),
  (2, 'user', 9, 1, NOW(), NOW());

-- step 3: siapa pun dari divisi finance
INSERT INTO `approval_flow_step_approvers`
  (`approval_flow_step_id`, `approver_type`, `approver_division_id`, `is_active`, `created_at`, `updated_at`)
VALUES
  (3, 'division', 2, 1, NOW(), NOW());
*/

-- =============================================================
-- Contoh alur transaksi request -> approval -> masuk project
-- =============================================================

-- 1) Public submit request (tanpa login)
/*
INSERT INTO `project_requests` (
  `request_no`, `project_name`, `project_description`, `business_goal`,
  `request_division_id`, `requested_ticket_prefix`,
  `requester_name`, `requester_email`, `requester_phone`, `requester_employee_id`,
  `approval_flow_id`, `current_step_order`, `status`, `created_at`, `updated_at`
)
VALUES (
  'PR-20260518-0001',
  'Portal Vendor',
  'Portal onboarding vendor',
  'Mempercepat proses procurement',
  2,
  'PVD',
  'Budi Santoso',
  'budi@company.com',
  '08123456789',
  'EMP-0092',
  1,
  1,
  'pending',
  NOW(),
  NOW()
);

-- snapshot step flow ke request (agar histori stabil walau flow master diubah)
INSERT INTO `project_request_step_states`
  (`project_request_id`, `approval_flow_step_id`, `step_order`, `step_name`, `approval_rule`, `status`, `created_at`, `updated_at`)
SELECT
  1,
  s.id,
  s.step_order,
  s.step_name,
  s.approval_rule,
  'pending',
  NOW(),
  NOW()
FROM `approval_flow_steps` s
WHERE s.approval_flow_id = 1
  AND s.is_active = 1
ORDER BY s.step_order;
*/

-- 2) Approver login lalu approve step tertentu
/*
INSERT INTO `project_request_approvals`
  (`project_request_id`, `approval_flow_step_id`, `approver_user_id`, `decision`, `note`, `created_at`)
VALUES
  (1, 1, 15, 'approved', 'OK lanjut', NOW());

-- jika rule step = any, 1 approval cukup:
UPDATE `project_request_step_states`
SET `status` = 'approved', `decided_by` = 15, `decided_at` = NOW(), `updated_at` = NOW()
WHERE `project_request_id` = 1
  AND `approval_flow_step_id` = 1
  AND `status` = 'pending';

-- pindah ke step berikutnya jika tidak ada step pending di order saat ini
UPDATE `project_requests`
SET `current_step_order` = `current_step_order` + 1, `updated_at` = NOW()
WHERE `id` = 1
  AND `status` = 'pending';
*/

-- 3) Rejected di step mana pun -> request selesai rejected
/*
INSERT INTO `project_request_approvals`
  (`project_request_id`, `approval_flow_step_id`, `approver_user_id`, `decision`, `note`, `created_at`)
VALUES
  (1, 2, 9, 'rejected', 'Belum siap dari sisi infrastruktur', NOW());

UPDATE `project_request_step_states`
SET `status` = 'rejected', `decided_by` = 9, `decided_at` = NOW(), `updated_at` = NOW()
WHERE `project_request_id` = 1
  AND `approval_flow_step_id` = 2
  AND `status` = 'pending';

UPDATE `project_requests`
SET
  `status` = 'rejected',
  `final_decided_by` = 9,
  `final_decided_at` = NOW(),
  `rejection_reason` = 'Belum siap dari sisi infrastruktur',
  `updated_at` = NOW()
WHERE `id` = 1
  AND `status` = 'pending';
*/

-- 4) Jika semua step approved -> sync ke projects status "Request Received"
/*
START TRANSACTION;

UPDATE `project_requests` pr
SET
  pr.`status` = 'approved',
  pr.`final_decided_by` = 3,
  pr.`final_decided_at` = NOW(),
  pr.`updated_at` = NOW()
WHERE pr.`id` = 1
  AND pr.`status` = 'pending'
  AND NOT EXISTS (
    SELECT 1
    FROM `project_request_step_states` s
    WHERE s.project_request_id = pr.id
      AND s.status IN ('pending', 'rejected')
  );

INSERT INTO `projects`
  (`name`, `description`, `owner_id`, `developer_id`, `status_id`, `priority_id`, `ticket_prefix`, `status_type`, `type`, `created_at`, `updated_at`)
SELECT
  pr.project_name,
  pr.project_description,
  pr.final_decided_by AS owner_id,
  pr.final_decided_by AS developer_id,
  ps.id AS status_id,
  pp.id AS priority_id,
  pr.requested_ticket_prefix,
  'default',
  'kanban',
  NOW(),
  NOW()
FROM `project_requests` pr
JOIN `project_statuses` ps
  ON LOWER(TRIM(ps.name)) = 'request received'
  AND ps.deleted_at IS NULL
JOIN `project_priorities` pp
  ON pp.is_default = 1
  AND pp.deleted_at IS NULL
LEFT JOIN `projects` p_exist
  ON p_exist.ticket_prefix = pr.requested_ticket_prefix
  AND p_exist.deleted_at IS NULL
WHERE pr.id = 1
  AND pr.status = 'approved'
  AND pr.project_id IS NULL
  AND p_exist.id IS NULL
LIMIT 1;

UPDATE `project_requests` pr
JOIN `projects` p
  ON p.ticket_prefix = pr.requested_ticket_prefix
  AND p.deleted_at IS NULL
SET
  pr.`project_id` = p.id,
  pr.`status` = 'synced_to_project',
  pr.`updated_at` = NOW()
WHERE pr.`id` = 1
  AND pr.`project_id` IS NULL
  AND pr.`status` = 'approved';

COMMIT;
*/
