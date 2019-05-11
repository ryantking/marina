CREATE TABLE `organization` (
	`name` VARCHAR(255) PRIMARY KEY,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE `repository` (
	`name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
	INDEX (`name`),
	FOREIGN KEY (`org_name`) REFERENCES `organization` (`name`)
);

CREATE TABLE `image` (
	`digest` VARCHAR(255) PRIMARY KEY,
	`repo_name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	`manifest` JSON NOT NULL,
	`manifest_type` VARCHAR(255) NOT NULL,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
	FOREIGN KEY (`repo_name`) REFERENCES `repository` (`name`),
	FOREIGN KEY (`org_name`) REFERENCES `organization` (`name`)
);

CREATE TABLE `tag` (
	`name` VARCHAR(255) NOT NULL,
	`repo_name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	`image_digest` VARCHAR(255) NOT NULL,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
	FOREIGN KEY (`image_digest`) REFERENCES `image` (`digest`)
);

CREATE TABLE `layer` (
	`digest` VARCHAR(255) PRIMARY KEY,
	`repo_name` VARCHAR(255) NOT NULL,
	`org_name` VARCHAR(255) NOT NULL,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW(),
	FOREIGN KEY (`repo_name`) REFERENCES `repository` (`name`),
	FOREIGN KEY (`org_name`) REFERENCES `organization` (`name`)
);

CREATE TABLE `upload` (
	`uuid` CHAR(36) PRIMARY KEY,
	`done` BOOLEAN NOT NULL DEFAULT FALSE,
	`_last_updated` TIMESTAMP NOT NULL DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE `upload_chunk` (
	`uuid` CHAR(36),
	`range_start` BIGINT,
	`range_end` BIGINT,
	PRIMARY KEY(`uuid`, `range_start`, `range_end`)
);

INSERT INTO `organization` (`name`)
VALUES ('library');
