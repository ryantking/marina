CREATE TABLE `organization` (
	`name` VARCHAR(255) PRIMARY KEY
);

CREATE TABLE `repository` (
	`name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	INDEX (`name`),
	FOREIGN KEY (`org_name`) REFERENCES `organization` (`name`)
);

CREATE TABLE `tag` (
	`name` VARCHAR(255) NOT NULL,
	`repo_name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	`manifest` JSON NOT NULL,
	`manifest_type` VARCHAR(255) NOT NULL,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
	FOREIGN KEY (`repo_name`) REFERENCES `repository` (`name`),
	FOREIGN KEY (`org_name`) REFERENCES `organization` (`name`)
);

CREATE TABLE `layer` (
	`digest` VARCHAR(1024) NOT NULL,
	`repo_name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	FOREIGN KEY (`repo_name`) REFERENCES `repository` (`name`),
	FOREIGN KEY (`org_name`) REFERENCES `organization` (`name`)
);

CREATE TABLE `upload` (
	`uuid` BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	`done` BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO `organization` (`name`)
VALUES ('library');
