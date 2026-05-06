-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 28, 2026 at 04:14 AM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `project-management-manna`
--

-- --------------------------------------------------------

--
-- Table structure for table `activities`
--

CREATE TABLE `activities` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` mediumtext NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `activities`
--

INSERT INTO `activities` (`id`, `name`, `description`, `created_at`, `updated_at`, `deleted_at`) VALUES
(1, 'Programming', 'Programming related activities', '2026-04-12 21:56:57', '2026-04-12 21:56:57', NULL),
(2, 'Testing', 'Testing related activities', '2026-04-12 21:56:57', '2026-04-12 21:56:57', NULL),
(3, 'Learning', 'Activities related to learning and training', '2026-04-12 21:56:57', '2026-04-12 21:56:57', NULL),
(4, 'Research', 'Activities related to research', '2026-04-12 21:56:57', '2026-04-12 21:56:57', NULL),
(5, 'Other', 'Other activities', '2026-04-12 21:56:57', '2026-04-12 21:56:57', NULL);

-- --------------------------------------------------------

--
-- Table structure for table `epics`
--

CREATE TABLE `epics` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `project_id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `starts_at` date NOT NULL,
  `ends_at` date NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `parent_id` bigint(20) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `epics`
--

INSERT INTO `epics` (`id`, `project_id`, `name`, `starts_at`, `ends_at`, `created_at`, `updated_at`, `deleted_at`, `parent_id`) VALUES
(1, 1, 'Requirement Gathering', '2026-04-13', '2026-06-18', '2026-04-12 23:00:19', '2026-04-27 00:35:10', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `failed_jobs`
--

CREATE TABLE `failed_jobs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `uuid` varchar(255) NOT NULL,
  `connection` text NOT NULL,
  `queue` text NOT NULL,
  `payload` longtext NOT NULL,
  `exception` longtext NOT NULL,
  `failed_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `jobs`
--

CREATE TABLE `jobs` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `queue` varchar(255) NOT NULL,
  `payload` longtext NOT NULL,
  `attempts` tinyint(3) UNSIGNED NOT NULL,
  `reserved_at` int(10) UNSIGNED DEFAULT NULL,
  `available_at` int(10) UNSIGNED NOT NULL,
  `created_at` int(10) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `jobs`
--

INSERT INTO `jobs` (`id`, `queue`, `payload`, `attempts`, `reserved_at`, `available_at`, `created_at`) VALUES
(1, 'default', '{\"uuid\":\"dcecbd1a-6e86-4551-acff-4d2c32c75154\",\"displayName\":\"App\\\\Notifications\\\\TicketCreated\",\"job\":\"Illuminate\\\\Queue\\\\CallQueuedHandler@call\",\"maxTries\":null,\"maxExceptions\":null,\"failOnTimeout\":false,\"backoff\":null,\"timeout\":null,\"retryUntil\":null,\"data\":{\"commandName\":\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\",\"command\":\"O:48:\\\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\\\":3:{s:11:\\\"notifiables\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:15:\\\"App\\\\Models\\\\User\\\";s:2:\\\"id\\\";a:1:{i:0;i:1;}s:9:\\\"relations\\\";a:0:{}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:12:\\\"notification\\\";O:31:\\\"App\\\\Notifications\\\\TicketCreated\\\":2:{s:39:\\\"\\u0000App\\\\Notifications\\\\TicketCreated\\u0000ticket\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:17:\\\"App\\\\Models\\\\Ticket\\\";s:2:\\\"id\\\";i:1;s:9:\\\"relations\\\";a:4:{i:0;s:7:\\\"project\\\";i:1;s:13:\\\"project.users\\\";i:2;s:5:\\\"owner\\\";i:3;s:11:\\\"responsible\\\";}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:2:\\\"id\\\";s:36:\\\"be59fb0f-5a8d-4e30-a664-be211089e1c3\\\";}s:8:\\\"channels\\\";a:1:{i:0;s:4:\\\"mail\\\";}}\"}}', 0, NULL, 1776060122, 1776060122),
(2, 'default', '{\"uuid\":\"f9faf1a9-5aba-4f2c-96fd-ed63b5e20d69\",\"displayName\":\"App\\\\Notifications\\\\TicketCreated\",\"job\":\"Illuminate\\\\Queue\\\\CallQueuedHandler@call\",\"maxTries\":null,\"maxExceptions\":null,\"failOnTimeout\":false,\"backoff\":null,\"timeout\":null,\"retryUntil\":null,\"data\":{\"commandName\":\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\",\"command\":\"O:48:\\\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\\\":3:{s:11:\\\"notifiables\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:15:\\\"App\\\\Models\\\\User\\\";s:2:\\\"id\\\";a:1:{i:0;i:1;}s:9:\\\"relations\\\";a:0:{}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:12:\\\"notification\\\";O:31:\\\"App\\\\Notifications\\\\TicketCreated\\\":2:{s:39:\\\"\\u0000App\\\\Notifications\\\\TicketCreated\\u0000ticket\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:17:\\\"App\\\\Models\\\\Ticket\\\";s:2:\\\"id\\\";i:1;s:9:\\\"relations\\\";a:4:{i:0;s:7:\\\"project\\\";i:1;s:13:\\\"project.users\\\";i:2;s:5:\\\"owner\\\";i:3;s:11:\\\"responsible\\\";}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:2:\\\"id\\\";s:36:\\\"be59fb0f-5a8d-4e30-a664-be211089e1c3\\\";}s:8:\\\"channels\\\";a:1:{i:0;s:8:\\\"database\\\";}}\"}}', 0, NULL, 1776060122, 1776060122),
(3, 'default', '{\"uuid\":\"fd754c8e-18b4-49c4-b0bd-3b3ee14196df\",\"displayName\":\"App\\\\Notifications\\\\TicketCreated\",\"job\":\"Illuminate\\\\Queue\\\\CallQueuedHandler@call\",\"maxTries\":null,\"maxExceptions\":null,\"failOnTimeout\":false,\"backoff\":null,\"timeout\":null,\"retryUntil\":null,\"data\":{\"commandName\":\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\",\"command\":\"O:48:\\\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\\\":3:{s:11:\\\"notifiables\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:15:\\\"App\\\\Models\\\\User\\\";s:2:\\\"id\\\";a:1:{i:0;i:1;}s:9:\\\"relations\\\";a:0:{}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:12:\\\"notification\\\";O:31:\\\"App\\\\Notifications\\\\TicketCreated\\\":2:{s:39:\\\"\\u0000App\\\\Notifications\\\\TicketCreated\\u0000ticket\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:17:\\\"App\\\\Models\\\\Ticket\\\";s:2:\\\"id\\\";i:2;s:9:\\\"relations\\\";a:4:{i:0;s:7:\\\"project\\\";i:1;s:13:\\\"project.users\\\";i:2;s:5:\\\"owner\\\";i:3;s:11:\\\"responsible\\\";}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:2:\\\"id\\\";s:36:\\\"e538fc4a-de26-429a-a00e-a05a3b067d1c\\\";}s:8:\\\"channels\\\";a:1:{i:0;s:4:\\\"mail\\\";}}\"}}', 0, NULL, 1777274587, 1777274587),
(4, 'default', '{\"uuid\":\"ef0234c6-6789-4fbf-91df-a28c432f4582\",\"displayName\":\"App\\\\Notifications\\\\TicketCreated\",\"job\":\"Illuminate\\\\Queue\\\\CallQueuedHandler@call\",\"maxTries\":null,\"maxExceptions\":null,\"failOnTimeout\":false,\"backoff\":null,\"timeout\":null,\"retryUntil\":null,\"data\":{\"commandName\":\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\",\"command\":\"O:48:\\\"Illuminate\\\\Notifications\\\\SendQueuedNotifications\\\":3:{s:11:\\\"notifiables\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:15:\\\"App\\\\Models\\\\User\\\";s:2:\\\"id\\\";a:1:{i:0;i:1;}s:9:\\\"relations\\\";a:0:{}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:12:\\\"notification\\\";O:31:\\\"App\\\\Notifications\\\\TicketCreated\\\":2:{s:39:\\\"\\u0000App\\\\Notifications\\\\TicketCreated\\u0000ticket\\\";O:45:\\\"Illuminate\\\\Contracts\\\\Database\\\\ModelIdentifier\\\":5:{s:5:\\\"class\\\";s:17:\\\"App\\\\Models\\\\Ticket\\\";s:2:\\\"id\\\";i:2;s:9:\\\"relations\\\";a:4:{i:0;s:7:\\\"project\\\";i:1;s:13:\\\"project.users\\\";i:2;s:5:\\\"owner\\\";i:3;s:11:\\\"responsible\\\";}s:10:\\\"connection\\\";s:5:\\\"mysql\\\";s:15:\\\"collectionClass\\\";N;}s:2:\\\"id\\\";s:36:\\\"e538fc4a-de26-429a-a00e-a05a3b067d1c\\\";}s:8:\\\"channels\\\";a:1:{i:0;s:8:\\\"database\\\";}}\"}}', 0, NULL, 1777274587, 1777274587);

-- --------------------------------------------------------

--
-- Table structure for table `media`
--

CREATE TABLE `media` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `model_type` varchar(255) NOT NULL,
  `model_id` bigint(20) UNSIGNED NOT NULL,
  `uuid` char(36) DEFAULT NULL,
  `collection_name` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `file_name` varchar(255) NOT NULL,
  `mime_type` varchar(255) DEFAULT NULL,
  `disk` varchar(255) NOT NULL,
  `conversions_disk` varchar(255) DEFAULT NULL,
  `size` bigint(20) UNSIGNED NOT NULL,
  `manipulations` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`manipulations`)),
  `custom_properties` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`custom_properties`)),
  `generated_conversions` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`generated_conversions`)),
  `responsive_images` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`responsive_images`)),
  `order_column` int(10) UNSIGNED DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `media`
--

INSERT INTO `media` (`id`, `model_type`, `model_id`, `uuid`, `collection_name`, `name`, `file_name`, `mime_type`, `disk`, `conversions_disk`, `size`, `manipulations`, `custom_properties`, `generated_conversions`, `responsive_images`, `order_column`, `created_at`, `updated_at`) VALUES
(1, 'App\\Models\\Project', 1, '8cf990ba-8869-4013-972d-91b577e197c1', 'default', 'Screenshot 2026-04-13 120546', 'dFw0I4UnMJdA2Z8ohCyD4mjNc83L7A-metaU2NyZWVuc2hvdCAyMDI2LTA0LTEzIDEyMDU0Ni5wbmc=-.png', 'image/png', 'public', 'public', 71019, '[]', '[]', '[]', '[]', 1, '2026-04-12 22:13:14', '2026-04-12 22:13:14');

-- --------------------------------------------------------

--
-- Table structure for table `migrations`
--

CREATE TABLE `migrations` (
  `id` int(10) UNSIGNED NOT NULL,
  `migration` varchar(255) NOT NULL,
  `batch` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `migrations`
--

INSERT INTO `migrations` (`id`, `migration`, `batch`) VALUES
(1, '2014_10_12_000000_create_users_table', 1),
(2, '2014_10_12_100000_create_password_resets_table', 1),
(3, '2019_08_19_000000_create_failed_jobs_table', 1),
(4, '2019_12_14_000001_create_personal_access_tokens_table', 1),
(5, '2022_11_02_111430_add_two_factor_columns_to_table', 1),
(6, '2022_11_02_113007_create_permission_tables', 1),
(7, '2022_11_02_124027_create_project_statuses_table', 1),
(8, '2022_11_02_124028_create_projects_table', 1),
(9, '2022_11_02_131753_create_project_users_table', 1),
(10, '2022_11_02_134510_create_media_table', 1),
(11, '2022_11_02_152359_create_project_favorites_table', 1),
(12, '2022_11_02_193241_create_ticket_statuses_table', 1),
(13, '2022_11_02_193242_create_tickets_table', 1),
(14, '2022_11_06_155109_add_tickets_prefix_to_projects', 1),
(15, '2022_11_06_163226_add_code_to_tickets', 1),
(16, '2022_11_06_164004_create_ticket_types_table', 1),
(17, '2022_11_06_165400_add_type_to_ticket', 1),
(18, '2022_11_06_173220_add_order_to_tickets', 1),
(19, '2022_11_06_184448_add_order_to_ticket_statuses', 1),
(20, '2022_11_06_193051_create_ticket_activities_table', 1),
(21, '2022_11_06_194000_create_ticket_priorities_table', 1),
(22, '2022_11_06_194728_add_priority_to_tickets', 1),
(23, '2022_11_06_203702_add_status_type_to_project', 1),
(24, '2022_11_06_204227_add_project_to_ticket_statuses', 1),
(25, '2022_11_07_064347_create_ticket_comments_table', 1),
(26, '2022_11_08_084509_create_ticket_subscribers_table', 1),
(27, '2022_11_08_144611_create_notifications_table', 1),
(28, '2022_11_08_150309_create_jobs_table', 1),
(29, '2022_11_08_163244_create_ticket_relations_table', 1),
(30, '2022_11_08_172846_create_settings_table', 1),
(31, '2022_11_08_173004_general_settings', 1),
(32, '2022_11_08_173852_create_general_settings', 1),
(33, '2022_11_09_085506_create_socialite_users_table', 1),
(34, '2022_11_09_085638_make_user_password_nullable', 1),
(35, '2022_11_09_110740_remove_unique_from_users', 1),
(36, '2022_11_09_110955_add_soft_deletes_to_users', 1),
(37, '2022_11_09_173852_add_social_login_to_general_settings', 1),
(38, '2022_11_10_193214_create_ticket_hours_table', 1),
(39, '2022_11_10_200608_add_estimation_to_tickets', 1),
(40, '2022_11_12_134201_add_creation_token_to_users', 1),
(41, '2022_11_12_142644_create_pending_user_emails_table', 1),
(42, '2022_11_12_173852_add_default_role_to_general_settings', 1),
(43, '2022_11_12_173852_add_login_form_oidc_enabled_flags_to_general_settings', 1),
(44, '2022_11_12_173852_add_site_language_to_general_settings', 1),
(45, '2022_12_15_100852_create_epics_table', 1),
(46, '2022_12_15_101035_add_epic_to_ticket', 1),
(47, '2022_12_16_133836_add_parent_to_epics', 1),
(48, '2022_12_27_082239_add_comment_to_ticket_hours', 1),
(49, '2023_01_05_182946_add_attachments_to_tickets', 1),
(50, '2023_01_09_113159_create_activities_table', 1),
(51, '2023_01_09_113847_add_activity_to_ticket_hours_table', 1),
(52, '2023_01_12_203211_remove_unique_constraint_from_users', 1),
(53, '2023_01_12_204221_drop_attachments', 1),
(54, '2023_01_15_201358_add_type_to_projects', 1),
(55, '2023_01_15_202225_create_sprints_table', 1),
(56, '2023_01_15_204606_add_sprint_to_tickets', 1),
(57, '2023_01_15_214849_add_epic_to_sprints', 1),
(58, '2023_01_16_085329_add_started_ended_at_to_sprints', 1),
(59, '2023_01_24_084637_update_users_for_oidc', 1),
(60, '2023_04_10_123922_add_unique_ticket_prefix_to_projects_table', 1);

-- --------------------------------------------------------

--
-- Table structure for table `model_has_permissions`
--

CREATE TABLE `model_has_permissions` (
  `permission_id` bigint(20) UNSIGNED NOT NULL,
  `model_type` varchar(255) NOT NULL,
  `model_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `model_has_roles`
--

CREATE TABLE `model_has_roles` (
  `role_id` bigint(20) UNSIGNED NOT NULL,
  `model_type` varchar(255) NOT NULL,
  `model_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `model_has_roles`
--

INSERT INTO `model_has_roles` (`role_id`, `model_type`, `model_id`) VALUES
(1, 'App\\Models\\User', 1);

-- --------------------------------------------------------

--
-- Table structure for table `notifications`
--

CREATE TABLE `notifications` (
  `id` char(36) NOT NULL,
  `type` varchar(255) NOT NULL,
  `notifiable_type` varchar(255) NOT NULL,
  `notifiable_id` bigint(20) UNSIGNED NOT NULL,
  `data` text NOT NULL,
  `read_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `password_resets`
--

CREATE TABLE `password_resets` (
  `email` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `pending_user_emails`
--

CREATE TABLE `pending_user_emails` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_type` varchar(255) NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `email` varchar(255) NOT NULL,
  `token` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `permissions`
--

CREATE TABLE `permissions` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `guard_name` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `permissions`
--

INSERT INTO `permissions` (`id`, `name`, `guard_name`, `created_at`, `updated_at`) VALUES
(1, 'List permissions', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(2, 'View permission', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(3, 'Create permission', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(4, 'Update permission', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(5, 'Delete permission', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(6, 'List projects', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(7, 'View project', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(8, 'Create project', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(9, 'Update project', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(10, 'Delete project', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(11, 'List project statuses', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(12, 'View project status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(13, 'Create project status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(14, 'Update project status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(15, 'Delete project status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(16, 'List roles', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(17, 'View role', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(18, 'Create role', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(19, 'Update role', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(20, 'Delete role', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(21, 'List tickets', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(22, 'View ticket', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(23, 'Create ticket', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(24, 'Update ticket', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(25, 'Delete ticket', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(26, 'List ticket priorities', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(27, 'View ticket priority', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(28, 'Create ticket priority', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(29, 'Update ticket priority', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(30, 'Delete ticket priority', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(31, 'List ticket statuses', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(32, 'View ticket status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(33, 'Create ticket status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(34, 'Update ticket status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(35, 'Delete ticket status', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(36, 'List ticket types', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(37, 'View ticket type', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(38, 'Create ticket type', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(39, 'Update ticket type', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(40, 'Delete ticket type', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(41, 'List users', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(42, 'View user', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(43, 'Create user', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(44, 'Update user', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(45, 'Delete user', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(46, 'List activities', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(47, 'View activity', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(48, 'Create activity', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(49, 'Update activity', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(50, 'Delete activity', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(51, 'List sprints', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(52, 'View sprint', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(53, 'Create sprint', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(54, 'Update sprint', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(55, 'Delete sprint', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(56, 'Manage general settings', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(57, 'Import from Jira', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(58, 'List timesheet data', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(59, 'View timesheet dashboard', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57');

-- --------------------------------------------------------

--
-- Table structure for table `personal_access_tokens`
--

CREATE TABLE `personal_access_tokens` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `tokenable_type` varchar(255) NOT NULL,
  `tokenable_id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `token` varchar(64) NOT NULL,
  `abilities` text DEFAULT NULL,
  `last_used_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `projects`
--

CREATE TABLE `projects` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` longtext DEFAULT NULL,
  `owner_id` bigint(20) UNSIGNED NOT NULL,
  `status_id` bigint(20) UNSIGNED NOT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `ticket_prefix` varchar(255) NOT NULL,
  `status_type` varchar(255) NOT NULL DEFAULT 'default',
  `type` varchar(255) NOT NULL DEFAULT 'kanban'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `projects`
--

INSERT INTO `projects` (`id`, `name`, `description`, `owner_id`, `status_id`, `deleted_at`, `created_at`, `updated_at`, `ticket_prefix`, `status_type`, `type`) VALUES
(1, 'Customer Care Manna Kampus', '<p>Manajemen Pengaduan<br>Kelola, prioritaskan, dan lacak keluhan customer dengan mudah. Tetapkan prioritas dan SLA untuk memastikan setiap masalah ditangani tepat waktu.</p>', 1, 1, NULL, '2026-04-12 22:13:14', '2026-04-12 22:13:14', 'SE', 'default', 'kanban');

-- --------------------------------------------------------

--
-- Table structure for table `project_favorites`
--

CREATE TABLE `project_favorites` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `project_id` bigint(20) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `project_statuses`
--

CREATE TABLE `project_statuses` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL DEFAULT '#cecece',
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `project_statuses`
--

INSERT INTO `project_statuses` (`id`, `name`, `color`, `is_default`, `deleted_at`, `created_at`, `updated_at`) VALUES
(1, 'Request Received', '#3b82f6', 1, NULL, '2026-04-12 22:07:26', '2026-04-12 22:08:46'),
(2, 'In Progress', '#f59e0b', 0, NULL, '2026-04-12 22:07:26', '2026-04-12 22:08:46'),
(3, 'Testing', '#8b5cf6', 0, NULL, '2026-04-12 22:07:26', '2026-04-12 22:08:46'),
(4, 'Implementation', '#f97316', 0, NULL, '2026-04-12 22:07:26', '2026-04-12 22:08:46'),
(5, 'Done', '#22c55e', 0, NULL, '2026-04-12 22:07:26', '2026-04-12 22:08:46');

-- --------------------------------------------------------

--
-- Table structure for table `project_users`
--

CREATE TABLE `project_users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `project_id` bigint(20) UNSIGNED NOT NULL,
  `role` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `project_users`
--

INSERT INTO `project_users` (`id`, `user_id`, `project_id`, `role`, `created_at`, `updated_at`) VALUES
(1, 1, 1, 'administrator', NULL, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `roles`
--

CREATE TABLE `roles` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `guard_name` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `roles`
--

INSERT INTO `roles` (`id`, `name`, `guard_name`, `created_at`, `updated_at`) VALUES
(1, 'Default role', 'web', '2026-04-12 21:56:57', '2026-04-12 21:56:57');

-- --------------------------------------------------------

--
-- Table structure for table `role_has_permissions`
--

CREATE TABLE `role_has_permissions` (
  `permission_id` bigint(20) UNSIGNED NOT NULL,
  `role_id` bigint(20) UNSIGNED NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `role_has_permissions`
--

INSERT INTO `role_has_permissions` (`permission_id`, `role_id`) VALUES
(1, 1),
(2, 1),
(3, 1),
(4, 1),
(5, 1),
(6, 1),
(7, 1),
(8, 1),
(9, 1),
(10, 1),
(11, 1),
(12, 1),
(13, 1),
(14, 1),
(15, 1),
(16, 1),
(17, 1),
(18, 1),
(19, 1),
(20, 1),
(21, 1),
(22, 1),
(23, 1),
(24, 1),
(25, 1),
(26, 1),
(27, 1),
(28, 1),
(29, 1),
(30, 1),
(31, 1),
(32, 1),
(33, 1),
(34, 1),
(35, 1),
(36, 1),
(37, 1),
(38, 1),
(39, 1),
(40, 1),
(41, 1),
(42, 1),
(43, 1),
(44, 1),
(45, 1),
(46, 1),
(47, 1),
(48, 1),
(49, 1),
(50, 1),
(51, 1),
(52, 1),
(53, 1),
(54, 1),
(55, 1),
(56, 1),
(57, 1),
(58, 1),
(59, 1);

-- --------------------------------------------------------

--
-- Table structure for table `settings`
--

CREATE TABLE `settings` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `group` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `locked` tinyint(1) NOT NULL DEFAULT 0,
  `payload` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL CHECK (json_valid(`payload`)),
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `settings`
--

INSERT INTO `settings` (`id`, `group`, `name`, `locked`, `payload`, `created_at`, `updated_at`) VALUES
(1, 'general', 'site_name', 0, '\"PM-MannaKampus\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(2, 'general', 'site_logo', 0, '\"A8281DLusKYNropcSaLt6nRHUyNieD-metaTG9nb19NYW5uYV9LYW1wdXMucG5n-.png\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(3, 'general', 'enable_registration', 0, 'true', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(4, 'general', 'enable_social_login', 0, '\"1\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(5, 'general', 'default_role', 0, '\"1\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(6, 'general', 'enable_login_form', 0, '\"1\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(7, 'general', 'enable_oidc_login', 0, '\"1\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58'),
(8, 'general', 'site_language', 0, '\"id\"', '2026-04-12 21:56:53', '2026-04-12 22:19:58');

-- --------------------------------------------------------

--
-- Table structure for table `socialite_users`
--

CREATE TABLE `socialite_users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `provider` varchar(255) NOT NULL,
  `provider_id` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `sprints`
--

CREATE TABLE `sprints` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `starts_at` date NOT NULL,
  `ends_at` date NOT NULL,
  `description` longtext DEFAULT NULL,
  `project_id` bigint(20) UNSIGNED NOT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `epic_id` bigint(20) UNSIGNED DEFAULT NULL,
  `started_at` datetime DEFAULT NULL,
  `ended_at` datetime DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `tickets`
--

CREATE TABLE `tickets` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `owner_id` bigint(20) UNSIGNED NOT NULL,
  `responsible_id` bigint(20) UNSIGNED DEFAULT NULL,
  `status_id` bigint(20) UNSIGNED NOT NULL,
  `project_id` bigint(20) UNSIGNED NOT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `code` varchar(255) NOT NULL,
  `type_id` bigint(20) UNSIGNED NOT NULL,
  `order` int(11) NOT NULL DEFAULT 0,
  `priority_id` bigint(20) UNSIGNED NOT NULL,
  `estimation` double(8,2) DEFAULT NULL,
  `starts_at` date DEFAULT NULL,
  `ends_at` date DEFAULT NULL,
  `epic_id` bigint(20) UNSIGNED DEFAULT NULL,
  `sprint_id` bigint(20) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `tickets`
--

INSERT INTO `tickets` (`id`, `name`, `content`, `owner_id`, `responsible_id`, `status_id`, `project_id`, `deleted_at`, `created_at`, `updated_at`, `code`, `type_id`, `order`, `priority_id`, `estimation`, `starts_at`, `ends_at`, `epic_id`, `sprint_id`) VALUES
(1, 'Meeting', '<p>Bahas Alur APP Request</p>', 1, 1, 1, 1, NULL, '2026-04-12 23:02:02', '2026-04-27 00:45:40', 'SE-1', 1, 0, 2, 8.00, '2026-04-14', '2026-04-15', 1, NULL),
(2, 'Buar Dockumentasi requirement Aplikasi', '<p>Buar Dockumentasi requirement Aplikasi</p>', 1, 1, 1, 1, NULL, '2026-04-27 00:23:07', '2026-04-27 00:24:38', 'SE-2', 1, 1, 2, 120.00, '2026-04-16', '2026-04-30', 1, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ticket_activities`
--

CREATE TABLE `ticket_activities` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `old_status_id` bigint(20) UNSIGNED NOT NULL,
  `new_status_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `ticket_attachments`
--

CREATE TABLE `ticket_attachments` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED DEFAULT NULL,
  `original_name` varchar(255) NOT NULL,
  `file_name` varchar(255) NOT NULL,
  `file_path` varchar(500) NOT NULL,
  `file_size` bigint(20) UNSIGNED NOT NULL DEFAULT 0,
  `mime_type` varchar(255) DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `ticket_comments`
--

CREATE TABLE `ticket_comments` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `content` longtext NOT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `ticket_hours`
--

CREATE TABLE `ticket_hours` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `value` double(8,2) NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `comment` longtext DEFAULT NULL,
  `activity_id` bigint(20) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `ticket_hours`
--

INSERT INTO `ticket_hours` (`id`, `ticket_id`, `user_id`, `value`, `created_at`, `updated_at`, `comment`, `activity_id`) VALUES
(1, 1, 1, 12.00, '2026-04-27 00:36:48', '2026-04-27 00:36:48', NULL, 4),
(2, 1, 1, 11.00, '2026-04-27 00:47:51', '2026-04-27 00:47:51', NULL, 5),
(3, 2, 1, 1.00, '2026-04-27 00:48:25', '2026-04-27 00:48:25', NULL, 5);

-- --------------------------------------------------------

--
-- Table structure for table `ticket_priorities`
--

CREATE TABLE `ticket_priorities` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL DEFAULT '#cecece',
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `ticket_priorities`
--

INSERT INTO `ticket_priorities` (`id`, `name`, `color`, `is_default`, `deleted_at`, `created_at`, `updated_at`) VALUES
(1, 'Low', '#008000', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(2, 'Normal', '#CECECE', 1, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(3, 'High', '#ff0000', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57');

-- --------------------------------------------------------

--
-- Table structure for table `ticket_relations`
--

CREATE TABLE `ticket_relations` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `relation_id` bigint(20) UNSIGNED NOT NULL,
  `type` varchar(255) NOT NULL,
  `sort` int(11) NOT NULL DEFAULT 1,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `ticket_statuses`
--

CREATE TABLE `ticket_statuses` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL DEFAULT '#cecece',
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `order` int(11) NOT NULL DEFAULT 1,
  `project_id` bigint(20) UNSIGNED DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `ticket_statuses`
--

INSERT INTO `ticket_statuses` (`id`, `name`, `color`, `is_default`, `deleted_at`, `created_at`, `updated_at`, `order`, `project_id`) VALUES
(1, 'Todo', '#cecece', 1, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57', 1, NULL),
(2, 'In progress', '#ff7f00', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57', 2, NULL),
(3, 'Done', '#008000', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57', 3, NULL),
(4, 'Archived', '#ff0000', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57', 4, NULL);

-- --------------------------------------------------------

--
-- Table structure for table `ticket_subscribers`
--

CREATE TABLE `ticket_subscribers` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `user_id` bigint(20) UNSIGNED NOT NULL,
  `ticket_id` bigint(20) UNSIGNED NOT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `ticket_types`
--

CREATE TABLE `ticket_types` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `icon` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL DEFAULT '#cecece',
  `is_default` tinyint(1) NOT NULL DEFAULT 0,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `ticket_types`
--

INSERT INTO `ticket_types` (`id`, `name`, `icon`, `color`, `is_default`, `deleted_at`, `created_at`, `updated_at`) VALUES
(1, 'Task', 'heroicon-o-check-circle', '#00FFFF', 1, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(2, 'Evolution', 'heroicon-o-clipboard-list', '#008000', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57'),
(3, 'Bug', 'heroicon-o-x', '#ff0000', 0, NULL, '2026-04-12 21:56:57', '2026-04-12 21:56:57');

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` bigint(20) UNSIGNED NOT NULL,
  `name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `email_verified_at` timestamp NULL DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  `two_factor_secret` text DEFAULT NULL,
  `two_factor_recovery_codes` text DEFAULT NULL,
  `two_factor_confirmed_at` timestamp NULL DEFAULT NULL,
  `remember_token` varchar(100) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  `creation_token` char(36) DEFAULT NULL,
  `type` varchar(255) NOT NULL DEFAULT 'db',
  `oidc_username` varchar(255) DEFAULT NULL,
  `oidc_sub` varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `name`, `email`, `email_verified_at`, `password`, `two_factor_secret`, `two_factor_recovery_codes`, `two_factor_confirmed_at`, `remember_token`, `created_at`, `updated_at`, `deleted_at`, `creation_token`, `type`, `oidc_username`, `oidc_sub`) VALUES
(1, 'Rifki Ahmad', 'admin@gmail.com', '2026-04-12 21:56:57', '$2y$10$iycFuTntX6Uc8FueGz9mb.Yr/1UOJ8foB3EVVnIpGz.wce47gsHQm', NULL, NULL, NULL, NULL, '2026-04-12 21:56:57', '2026-04-12 22:21:46', NULL, NULL, 'db', NULL, NULL),
(2, 'admin', 'asd@gmail.com', '2026-04-27 07:16:15', '$2y$10$iycFuTntX6Uc8FueGz9mb.Yr/1UOJ8foB3EVVnIpGz.wce47gsHQm', NULL, NULL, NULL, NULL, '2026-04-27 00:15:26', '2026-04-27 00:15:26', NULL, NULL, 'db', NULL, NULL);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `activities`
--
ALTER TABLE `activities`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `epics`
--
ALTER TABLE `epics`
  ADD PRIMARY KEY (`id`),
  ADD KEY `epics_project_id_foreign` (`project_id`),
  ADD KEY `epics_parent_id_foreign` (`parent_id`);

--
-- Indexes for table `failed_jobs`
--
ALTER TABLE `failed_jobs`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `failed_jobs_uuid_unique` (`uuid`);

--
-- Indexes for table `jobs`
--
ALTER TABLE `jobs`
  ADD PRIMARY KEY (`id`),
  ADD KEY `jobs_queue_index` (`queue`);

--
-- Indexes for table `media`
--
ALTER TABLE `media`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `media_uuid_unique` (`uuid`),
  ADD KEY `media_model_type_model_id_index` (`model_type`,`model_id`),
  ADD KEY `media_order_column_index` (`order_column`);

--
-- Indexes for table `migrations`
--
ALTER TABLE `migrations`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `model_has_permissions`
--
ALTER TABLE `model_has_permissions`
  ADD PRIMARY KEY (`permission_id`,`model_id`,`model_type`),
  ADD KEY `model_has_permissions_model_id_model_type_index` (`model_id`,`model_type`);

--
-- Indexes for table `model_has_roles`
--
ALTER TABLE `model_has_roles`
  ADD PRIMARY KEY (`role_id`,`model_id`,`model_type`),
  ADD KEY `model_has_roles_model_id_model_type_index` (`model_id`,`model_type`);

--
-- Indexes for table `notifications`
--
ALTER TABLE `notifications`
  ADD PRIMARY KEY (`id`),
  ADD KEY `notifications_notifiable_type_notifiable_id_index` (`notifiable_type`,`notifiable_id`);

--
-- Indexes for table `password_resets`
--
ALTER TABLE `password_resets`
  ADD KEY `password_resets_email_index` (`email`);

--
-- Indexes for table `pending_user_emails`
--
ALTER TABLE `pending_user_emails`
  ADD PRIMARY KEY (`id`),
  ADD KEY `pending_user_emails_user_type_user_id_index` (`user_type`,`user_id`),
  ADD KEY `pending_user_emails_email_index` (`email`);

--
-- Indexes for table `permissions`
--
ALTER TABLE `permissions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `permissions_name_guard_name_unique` (`name`,`guard_name`);

--
-- Indexes for table `personal_access_tokens`
--
ALTER TABLE `personal_access_tokens`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `personal_access_tokens_token_unique` (`token`),
  ADD KEY `personal_access_tokens_tokenable_type_tokenable_id_index` (`tokenable_type`,`tokenable_id`);

--
-- Indexes for table `projects`
--
ALTER TABLE `projects`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `projects_ticket_prefix_unique` (`ticket_prefix`),
  ADD KEY `projects_owner_id_foreign` (`owner_id`),
  ADD KEY `projects_status_id_foreign` (`status_id`);

--
-- Indexes for table `project_favorites`
--
ALTER TABLE `project_favorites`
  ADD PRIMARY KEY (`id`),
  ADD KEY `project_favorites_user_id_foreign` (`user_id`),
  ADD KEY `project_favorites_project_id_foreign` (`project_id`);

--
-- Indexes for table `project_statuses`
--
ALTER TABLE `project_statuses`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `project_users`
--
ALTER TABLE `project_users`
  ADD PRIMARY KEY (`id`),
  ADD KEY `project_users_user_id_foreign` (`user_id`),
  ADD KEY `project_users_project_id_foreign` (`project_id`);

--
-- Indexes for table `roles`
--
ALTER TABLE `roles`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `roles_name_guard_name_unique` (`name`,`guard_name`);

--
-- Indexes for table `role_has_permissions`
--
ALTER TABLE `role_has_permissions`
  ADD PRIMARY KEY (`permission_id`,`role_id`),
  ADD KEY `role_has_permissions_role_id_foreign` (`role_id`);

--
-- Indexes for table `settings`
--
ALTER TABLE `settings`
  ADD PRIMARY KEY (`id`),
  ADD KEY `settings_group_index` (`group`);

--
-- Indexes for table `socialite_users`
--
ALTER TABLE `socialite_users`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `socialite_users_provider_provider_id_unique` (`provider`,`provider_id`);

--
-- Indexes for table `sprints`
--
ALTER TABLE `sprints`
  ADD PRIMARY KEY (`id`),
  ADD KEY `sprints_project_id_foreign` (`project_id`),
  ADD KEY `sprints_epic_id_foreign` (`epic_id`);

--
-- Indexes for table `tickets`
--
ALTER TABLE `tickets`
  ADD PRIMARY KEY (`id`),
  ADD KEY `tickets_owner_id_foreign` (`owner_id`),
  ADD KEY `tickets_responsible_id_foreign` (`responsible_id`),
  ADD KEY `tickets_status_id_foreign` (`status_id`),
  ADD KEY `tickets_project_id_foreign` (`project_id`),
  ADD KEY `tickets_type_id_foreign` (`type_id`),
  ADD KEY `tickets_priority_id_foreign` (`priority_id`),
  ADD KEY `tickets_epic_id_foreign` (`epic_id`),
  ADD KEY `tickets_sprint_id_foreign` (`sprint_id`);

--
-- Indexes for table `ticket_activities`
--
ALTER TABLE `ticket_activities`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_activities_ticket_id_foreign` (`ticket_id`),
  ADD KEY `ticket_activities_old_status_id_foreign` (`old_status_id`),
  ADD KEY `ticket_activities_new_status_id_foreign` (`new_status_id`),
  ADD KEY `ticket_activities_user_id_foreign` (`user_id`);

--
-- Indexes for table `ticket_attachments`
--
ALTER TABLE `ticket_attachments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_attachments_ticket_id_foreign` (`ticket_id`),
  ADD KEY `ticket_attachments_user_id_foreign` (`user_id`);

--
-- Indexes for table `ticket_comments`
--
ALTER TABLE `ticket_comments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_comments_ticket_id_foreign` (`ticket_id`),
  ADD KEY `ticket_comments_user_id_foreign` (`user_id`);

--
-- Indexes for table `ticket_hours`
--
ALTER TABLE `ticket_hours`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_hours_ticket_id_foreign` (`ticket_id`),
  ADD KEY `ticket_hours_user_id_foreign` (`user_id`),
  ADD KEY `ticket_hours_activity_id_foreign` (`activity_id`);

--
-- Indexes for table `ticket_priorities`
--
ALTER TABLE `ticket_priorities`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `ticket_relations`
--
ALTER TABLE `ticket_relations`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_relations_ticket_id_foreign` (`ticket_id`),
  ADD KEY `ticket_relations_relation_id_foreign` (`relation_id`);

--
-- Indexes for table `ticket_statuses`
--
ALTER TABLE `ticket_statuses`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_statuses_project_id_foreign` (`project_id`);

--
-- Indexes for table `ticket_subscribers`
--
ALTER TABLE `ticket_subscribers`
  ADD PRIMARY KEY (`id`),
  ADD KEY `ticket_subscribers_user_id_foreign` (`user_id`),
  ADD KEY `ticket_subscribers_ticket_id_foreign` (`ticket_id`);

--
-- Indexes for table `ticket_types`
--
ALTER TABLE `ticket_types`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `activities`
--
ALTER TABLE `activities`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `epics`
--
ALTER TABLE `epics`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `failed_jobs`
--
ALTER TABLE `failed_jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `jobs`
--
ALTER TABLE `jobs`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `media`
--
ALTER TABLE `media`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `migrations`
--
ALTER TABLE `migrations`
  MODIFY `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=61;

--
-- AUTO_INCREMENT for table `pending_user_emails`
--
ALTER TABLE `pending_user_emails`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `permissions`
--
ALTER TABLE `permissions`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=60;

--
-- AUTO_INCREMENT for table `personal_access_tokens`
--
ALTER TABLE `personal_access_tokens`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `projects`
--
ALTER TABLE `projects`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `project_favorites`
--
ALTER TABLE `project_favorites`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `project_statuses`
--
ALTER TABLE `project_statuses`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `project_users`
--
ALTER TABLE `project_users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `roles`
--
ALTER TABLE `roles`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `settings`
--
ALTER TABLE `settings`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=9;

--
-- AUTO_INCREMENT for table `socialite_users`
--
ALTER TABLE `socialite_users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sprints`
--
ALTER TABLE `sprints`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `tickets`
--
ALTER TABLE `tickets`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- AUTO_INCREMENT for table `ticket_activities`
--
ALTER TABLE `ticket_activities`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticket_attachments`
--
ALTER TABLE `ticket_attachments`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticket_comments`
--
ALTER TABLE `ticket_comments`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticket_hours`
--
ALTER TABLE `ticket_hours`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `ticket_priorities`
--
ALTER TABLE `ticket_priorities`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `ticket_relations`
--
ALTER TABLE `ticket_relations`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticket_statuses`
--
ALTER TABLE `ticket_statuses`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `ticket_subscribers`
--
ALTER TABLE `ticket_subscribers`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `ticket_types`
--
ALTER TABLE `ticket_types`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `epics`
--
ALTER TABLE `epics`
  ADD CONSTRAINT `epics_parent_id_foreign` FOREIGN KEY (`parent_id`) REFERENCES `epics` (`id`),
  ADD CONSTRAINT `epics_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

--
-- Constraints for table `model_has_permissions`
--
ALTER TABLE `model_has_permissions`
  ADD CONSTRAINT `model_has_permissions_permission_id_foreign` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `model_has_roles`
--
ALTER TABLE `model_has_roles`
  ADD CONSTRAINT `model_has_roles_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `projects`
--
ALTER TABLE `projects`
  ADD CONSTRAINT `projects_owner_id_foreign` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`),
  ADD CONSTRAINT `projects_status_id_foreign` FOREIGN KEY (`status_id`) REFERENCES `project_statuses` (`id`);

--
-- Constraints for table `project_favorites`
--
ALTER TABLE `project_favorites`
  ADD CONSTRAINT `project_favorites_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`),
  ADD CONSTRAINT `project_favorites_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `project_users`
--
ALTER TABLE `project_users`
  ADD CONSTRAINT `project_users_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`),
  ADD CONSTRAINT `project_users_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `role_has_permissions`
--
ALTER TABLE `role_has_permissions`
  ADD CONSTRAINT `role_has_permissions_permission_id_foreign` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`) ON DELETE CASCADE,
  ADD CONSTRAINT `role_has_permissions_role_id_foreign` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE;

--
-- Constraints for table `sprints`
--
ALTER TABLE `sprints`
  ADD CONSTRAINT `sprints_epic_id_foreign` FOREIGN KEY (`epic_id`) REFERENCES `epics` (`id`),
  ADD CONSTRAINT `sprints_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

--
-- Constraints for table `tickets`
--
ALTER TABLE `tickets`
  ADD CONSTRAINT `tickets_epic_id_foreign` FOREIGN KEY (`epic_id`) REFERENCES `epics` (`id`),
  ADD CONSTRAINT `tickets_owner_id_foreign` FOREIGN KEY (`owner_id`) REFERENCES `users` (`id`),
  ADD CONSTRAINT `tickets_priority_id_foreign` FOREIGN KEY (`priority_id`) REFERENCES `ticket_priorities` (`id`),
  ADD CONSTRAINT `tickets_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`),
  ADD CONSTRAINT `tickets_responsible_id_foreign` FOREIGN KEY (`responsible_id`) REFERENCES `users` (`id`),
  ADD CONSTRAINT `tickets_sprint_id_foreign` FOREIGN KEY (`sprint_id`) REFERENCES `sprints` (`id`),
  ADD CONSTRAINT `tickets_status_id_foreign` FOREIGN KEY (`status_id`) REFERENCES `ticket_statuses` (`id`),
  ADD CONSTRAINT `tickets_type_id_foreign` FOREIGN KEY (`type_id`) REFERENCES `ticket_types` (`id`);

--
-- Constraints for table `ticket_activities`
--
ALTER TABLE `ticket_activities`
  ADD CONSTRAINT `ticket_activities_new_status_id_foreign` FOREIGN KEY (`new_status_id`) REFERENCES `ticket_statuses` (`id`),
  ADD CONSTRAINT `ticket_activities_old_status_id_foreign` FOREIGN KEY (`old_status_id`) REFERENCES `ticket_statuses` (`id`),
  ADD CONSTRAINT `ticket_activities_ticket_id_foreign` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`),
  ADD CONSTRAINT `ticket_activities_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `ticket_attachments`
--
ALTER TABLE `ticket_attachments`
  ADD CONSTRAINT `ticket_attachments_ticket_id_foreign` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`),
  ADD CONSTRAINT `ticket_attachments_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `ticket_comments`
--
ALTER TABLE `ticket_comments`
  ADD CONSTRAINT `ticket_comments_ticket_id_foreign` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`),
  ADD CONSTRAINT `ticket_comments_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `ticket_hours`
--
ALTER TABLE `ticket_hours`
  ADD CONSTRAINT `ticket_hours_activity_id_foreign` FOREIGN KEY (`activity_id`) REFERENCES `activities` (`id`),
  ADD CONSTRAINT `ticket_hours_ticket_id_foreign` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`),
  ADD CONSTRAINT `ticket_hours_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);

--
-- Constraints for table `ticket_relations`
--
ALTER TABLE `ticket_relations`
  ADD CONSTRAINT `ticket_relations_relation_id_foreign` FOREIGN KEY (`relation_id`) REFERENCES `tickets` (`id`),
  ADD CONSTRAINT `ticket_relations_ticket_id_foreign` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`);

--
-- Constraints for table `ticket_statuses`
--
ALTER TABLE `ticket_statuses`
  ADD CONSTRAINT `ticket_statuses_project_id_foreign` FOREIGN KEY (`project_id`) REFERENCES `projects` (`id`);

--
-- Constraints for table `ticket_subscribers`
--
ALTER TABLE `ticket_subscribers`
  ADD CONSTRAINT `ticket_subscribers_ticket_id_foreign` FOREIGN KEY (`ticket_id`) REFERENCES `tickets` (`id`),
  ADD CONSTRAINT `ticket_subscribers_user_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
