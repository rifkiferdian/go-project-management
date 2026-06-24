INSERT INTO `permissions` (`name`, `guard_name`, `created_at`, `updated_at`)
SELECT 'Copy ticket template', 'web', NOW(), NOW()
WHERE NOT EXISTS (
    SELECT 1
    FROM `permissions`
    WHERE `name` = 'Copy ticket template'
      AND `guard_name` = 'web'
);

-- Permission sengaja tidak diberikan otomatis ke role mana pun.
-- Atur akses melalui menu Roles.
