ALTER TABLE `approval_flows`
    ADD COLUMN `division_id` BIGINT(20) UNSIGNED NULL AFTER `id`,
    ADD KEY `approval_flows_division_id_foreign` (`division_id`),
    ADD CONSTRAINT `approval_flows_division_id_foreign`
        FOREIGN KEY (`division_id`) REFERENCES `divisions` (`id`);
