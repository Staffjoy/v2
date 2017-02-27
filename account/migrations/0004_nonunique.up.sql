DROP INDEX `ix_account_phonenumber` ON `account`;
DROP INDEX `ix_account_email` ON `account`;
UPDATE account SET email="" WHERE email IS NULL;
UPDATE account SET phonenumber="" WHERE phonenumber IS NULL;
ALTER TABLE `account` MODIFY COLUMN `phonenumber` VARCHAR(255) NOT NULL DEFAULT "";
ALTER TABLE `account` MODIFY COLUMN `email` VARCHAR(255) NOT NULL DEFAULT "";