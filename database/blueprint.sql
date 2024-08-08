CREATE TABLE `task` (
	 `id` char(36) PRIMARY KEY DEFAULT(UUID()),
	 `name` VARCHAR(50),
	 `description` VARCHAR(750)
);

INSERT INTO `task`(`name`, `description`) 
VALUES('Some name', 'Some description');
