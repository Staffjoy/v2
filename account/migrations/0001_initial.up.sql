CREATE TABLE `account` (
    `uuid` VARCHAR(255),
    `email` VARCHAR(255),
    `name` VARCHAR(255),
    `confirmed_and_active` TINYINT(1) DEFAULT 0,
    `member_since` DATETIME,
    `password_hash` BINARY(60),
    PRIMARY KEY (`uuid`),
    UNIQUE KEY `ix_account_email` (`email`)
) ENGINE=InnoDB;
