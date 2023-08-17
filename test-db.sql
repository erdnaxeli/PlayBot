/* File generated with "show create" on phpMyAdmin */
CREATE TABLE `playbot` (
    `type` varchar(15) NOT NULL,
    `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
    `sender_irc` varchar(99) DEFAULT NULL,
    `sender` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci DEFAULT NULL,
    `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
    `duration` int(6) DEFAULT NULL,
    `file` varchar(150) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
    `broken` tinyint(1) NOT NULL DEFAULT 0,
    `id` int(6) NOT NULL AUTO_INCREMENT,
    `channel` varchar(255) DEFAULT NULL,
    `playlist` tinyint(1) DEFAULT 0,
    `external_id` varchar(255) DEFAULT NULL,
    `eatshit` float DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `url` (`url`),
    KEY `chan` (`channel`),
    FULLTEXT KEY `sender` (`sender`, `title`)
) ENGINE = InnoDB AUTO_INCREMENT = 100 DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci COMMENT = 'table du bot irc';
CREATE TABLE `playbot_chan` (
    `id` int(10) NOT NULL AUTO_INCREMENT,
    `date` timestamp(6) NULL DEFAULT current_timestamp(),
    `sender_irc` varchar(99) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
    `content` int(10) NOT NULL,
    `chan` varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `content` (`content`),
    KEY `chan` (`chan`),
    CONSTRAINT `playbot_chan_ibfk_1` FOREIGN KEY (`content`) REFERENCES `playbot` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB AUTO_INCREMENT = 100 DEFAULT CHARSET = latin1 COLLATE = latin1_swedish_ci;
CREATE TABLE `playbot_codes` (
    `user` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
    `code` varchar(25) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
    `nick` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL,
    PRIMARY KEY (`user`),
    UNIQUE KEY `code` (`code`),
    UNIQUE KEY `nick` (`nick`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;
CREATE TABLE `playbot_config` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(255) NOT NULL,
    `sites` varchar(300) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 100 DEFAULT CHARSET = latin1 COLLATE = latin1_swedish_ci;
CREATE TABLE `playbot_fav` (
    `user` varchar(255) NOT NULL,
    `id` int(6) NOT NULL,
    `date` timestamp NULL DEFAULT current_timestamp(),
    PRIMARY KEY (`user`, `id`),
    KEY `id` (`id`),
    KEY `user` (`user`),
    CONSTRAINT `playbot_fav_ibfk_1` FOREIGN KEY (`user`) REFERENCES `playbot_codes` (`user`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = latin1 COLLATE = latin1_swedish_ci;
CREATE TABLE `playbot_later` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `content` int(11) NOT NULL,
    `nick` varchar(255) NOT NULL,
    `date` int(10) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `content` (`content`)
) ENGINE = InnoDB AUTO_INCREMENT = 100 DEFAULT CHARSET = latin1 COLLATE = latin1_swedish_ci;
CREATE TABLE `playbot_playlist_content_association` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `playlist_id` int(11) DEFAULT NULL,
    `content_id` int(11) DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `playlist_id` (`playlist_id`, `content_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 100 DEFAULT CHARSET = latin1 COLLATE = latin1_swedish_ci;
CREATE TABLE `playbot_tags` (
    `id` int(11) NOT NULL,
    `tag` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
    PRIMARY KEY (`id`, `tag`),
    KEY `tag` (`tag`),
    KEY `id` (`id`),
    CONSTRAINT `playbot_tags_ibfk_1` FOREIGN KEY (`id`) REFERENCES `playbot` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB DEFAULT CHARSET = utf8 COLLATE = utf8_unicode_ci;
