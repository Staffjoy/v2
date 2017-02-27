ALTER TABLE `account` ADD `phonenumber` VARCHAR(15);
ALTER TABLE `account` ADD `password_salt` BINARY(60) NOT NULL DEFAULT "\0";
ALTER TABLE `account` ADD `support` TINYINT(1) NOT NULL DEFAULT 0;
CREATE UNIQUE INDEX `ix_account_phonenumber` ON `account` (`phonenumber`);


ALTER TABLE `account` MODIFY `email` VARCHAR(255);
ALTER TABLE `account` MODIFY `name` VARCHAR(255) NOT NULL DEFAULT "";
ALTER TABLE `account` MODIFY `confirmed_and_active` TINYINT(0) NOT NULL DEFAULT 0;
ALTER TABLE `account` MODIFY `member_since` DATETIME NOT NULL;
ALTER TABLE `account` MODIFY `password_hash` BINARY(60) NOT NULL DEFAULT "";
