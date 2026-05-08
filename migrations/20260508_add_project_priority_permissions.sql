START TRANSACTION;

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'List project priorities', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'List project priorities' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'View project priority', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'View project priority' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Create project priority', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Create project priority' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Update project priority', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Update project priority' AND `guard_name` = 'web'
);

INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Delete project priority', 'web', NOW(), NOW()
WHERE NOT EXISTS (
  SELECT 1 FROM `permissions` WHERE `name` = 'Delete project priority' AND `guard_name` = 'web'
);

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT pp.id, rhp.role_id
FROM `permissions` pp
JOIN `permissions` tp ON tp.name = 'List ticket priorities' AND tp.guard_name = 'web'
JOIN `role_has_permissions` rhp ON rhp.permission_id = tp.id
LEFT JOIN `role_has_permissions` existing ON existing.permission_id = pp.id AND existing.role_id = rhp.role_id
WHERE pp.name = 'List project priorities'
  AND pp.guard_name = 'web'
  AND existing.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT pp.id, rhp.role_id
FROM `permissions` pp
JOIN `permissions` tp ON tp.name = 'View ticket priority' AND tp.guard_name = 'web'
JOIN `role_has_permissions` rhp ON rhp.permission_id = tp.id
LEFT JOIN `role_has_permissions` existing ON existing.permission_id = pp.id AND existing.role_id = rhp.role_id
WHERE pp.name = 'View project priority'
  AND pp.guard_name = 'web'
  AND existing.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT pp.id, rhp.role_id
FROM `permissions` pp
JOIN `permissions` tp ON tp.name = 'Create ticket priority' AND tp.guard_name = 'web'
JOIN `role_has_permissions` rhp ON rhp.permission_id = tp.id
LEFT JOIN `role_has_permissions` existing ON existing.permission_id = pp.id AND existing.role_id = rhp.role_id
WHERE pp.name = 'Create project priority'
  AND pp.guard_name = 'web'
  AND existing.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT pp.id, rhp.role_id
FROM `permissions` pp
JOIN `permissions` tp ON tp.name = 'Update ticket priority' AND tp.guard_name = 'web'
JOIN `role_has_permissions` rhp ON rhp.permission_id = tp.id
LEFT JOIN `role_has_permissions` existing ON existing.permission_id = pp.id AND existing.role_id = rhp.role_id
WHERE pp.name = 'Update project priority'
  AND pp.guard_name = 'web'
  AND existing.permission_id IS NULL;

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`)
SELECT pp.id, rhp.role_id
FROM `permissions` pp
JOIN `permissions` tp ON tp.name = 'Delete ticket priority' AND tp.guard_name = 'web'
JOIN `role_has_permissions` rhp ON rhp.permission_id = tp.id
LEFT JOIN `role_has_permissions` existing ON existing.permission_id = pp.id AND existing.role_id = rhp.role_id
WHERE pp.name = 'Delete project priority'
  AND pp.guard_name = 'web'
  AND existing.permission_id IS NULL;

COMMIT;
