USE url;

CREATE TABLE `redirect` (
	`token` varchar(14) collate utf8mb4_unicode_ci NOT NULL,
	`url` varchar(620) collate utf8mb4_unicode_ci NOT NULL,
	PRIMARY KEY (`token`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='Used for the URL shortener';

INSERT INTO `redirect` VALUES ('cogo', 'https://www.cogolabs.com/');
