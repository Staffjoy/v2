CREATE TABLE `company` (
    `uuid` VARCHAR(255),
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `archived` TINYINT(1) DEFAULT 0 NOT NULL,
    `default_timezone` VARCHAR(255) NOT NULL DEFAULT "",
    `default_day_week_starts` VARCHAR(20) NOT NULL DEFAULT "monday",
    PRIMARY KEY (`uuid`)
);

CREATE TABLE `team` (
    `uuid` VARCHAR(255) NOT NULL,
    `company_uuid` VARCHAR(255) NOT NULL DEFAULT "",
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `archived` TINYINT(1) NOT NULL DEFAULT 0,
    `timezone` VARCHAR(255) NOT NULL DEFAULT "",
    `day_week_starts` VARCHAR(20) NOT NULL DEFAULT "monday",
    `color` VARCHAR(10) NOT NULL DEFAULT "48B7AB",
    PRIMARY KEY (`uuid`), 
    KEY `ix_team_company_uuid` (`company_uuid`)
);

CREATE TABLE `shift` (
    `uuid` VARCHAR(255) NOT NULL,
    `team_uuid` VARCHAR(255) NOT NULL DEFAULT "",
    `job_uuid` VARCHAR(255) NOT NULL DEFAULT "",
    `user_uuid` VARCHAR(255) NOT NULL DEFAULT "",
    `published` TINYINT(1) NOT NULL DEFAULT 0,
    `start` DATETIME NOT NULL,
    `stop` DATETIME NOT NULL, 
    PRIMARY KEY (`uuid`), 
    KEY `ix_job_shift_uuid` (`job_uuid`),
    KEY `ix_job_user_uuid` (`user_uuid`)
);

CREATE TABLE `job` (
    `uuid` VARCHAR(255) NOT NULL,
    `team_uuid` VARCHAR(255) NOT NULL DEFAULT "",
    `name` VARCHAR(255) NOT NULL DEFAULT "",
    `archived` TINYINT(1) NOT NULL DEFAULT 0,
    `color` VARCHAR(10) NOT NULL DEFAULT "48B7AB",
    PRIMARY KEY (`uuid`), 
    KEY `ix_job_team_uuid` (`team_uuid`)
);

CREATE TABLE `directory` (
    `company_uuid` VARCHAR(255) NOT NULL,
    `user_uuid` VARCHAR(255) NOT NULL,
    `internal_id` VARCHAR(255) NOT NULL,
    KEY `ix_directory_company_uuid` (`company_uuid`),
    KEY `ix_directory_user_uuid` (`user_uuid`),
    KEY `ix_directory_internal_id` (`internal_id`)
);

CREATE TABLE `worker` (
    `team_uuid` VARCHAR(255) NOT NULL,
    `user_uuid` VARCHAR(255) NOT NULL,
    KEY `ix_team_team_uuid` (`team_uuid`),
    KEY `ix_team_user_uuid` (`user_uuid`)
);

CREATE TABLE `manager` (
    `team_uuid` VARCHAR(255) NOT NULL,
    `user_uuid` VARCHAR(255) NOT NULL,
    KEY `ix_manager_team_uuid` (`team_uuid`),
    KEY `ix_manager_user_uuid` (`user_uuid`)
);

CREATE TABLE `admin` (
    `company_uuid` VARCHAR(255) NOT NULL,
    `user_uuid` VARCHAR(255) NOT NULL,
    KEY `ix_admin_company_uuid` (`company_uuid`),
    KEY `ix_admin_user_uuid` (`user_uuid`)
);
