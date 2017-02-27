CREATE UNIQUE INDEX `ix_account_phonenumber` ON `account` (`phonenumber`);
CREATE UNIQUE INDEX `ix_account_email` ON `account` (`email`);
ALTER TABLE `account` UPDATE `phonenumber` VARCHAR(255) NULL;
ALTER TABLE `account` UPDATE `email` VARCHAR(255) NULL;