-- phpMyAdmin SQL Dump
-- version 5.2.0-dev
-- https://www.phpmyadmin.net/
--
-- Host: mysql.iiens.net
-- Generation Time: Aug 20, 2023 at 06:56 PM
-- Server version: 10.5.19-MariaDB-0+deb11u2
-- PHP Version: 7.4.14

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `assoce_nightiies`
--

-- --------------------------------------------------------

--
-- Table structure for table `playbot`
--

CREATE TABLE `playbot` (
  `type` varchar(15) NOT NULL,
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `sender_irc` varchar(99) DEFAULT NULL,
  `sender` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci DEFAULT NULL,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `duration` int(6) DEFAULT NULL,
  `file` varchar(150) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `broken` tinyint(1) NOT NULL DEFAULT 0,
  `id` int(6) NOT NULL,
  `channel` varchar(255) DEFAULT NULL,
  `playlist` tinyint(1) DEFAULT 0,
  `external_id` varchar(255) DEFAULT NULL,
  `eatshit` float DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='table du bot irc';

-- --------------------------------------------------------

--
-- Table structure for table `playbot_chan`
--

CREATE TABLE `playbot_chan` (
  `id` int(10) NOT NULL,
  `date` timestamp(6) NULL DEFAULT current_timestamp(6),
  `sender_irc` varchar(99) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  `content` int(10) NOT NULL,
  `chan` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `playbot_codes`
--

CREATE TABLE `playbot_codes` (
  `user` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `code` varchar(25) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL,
  `nick` varchar(255) CHARACTER SET latin1 COLLATE latin1_swedish_ci DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

-- --------------------------------------------------------

--
-- Table structure for table `playbot_config`
--

CREATE TABLE `playbot_config` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `sites` varchar(300) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `playbot_fav`
--

CREATE TABLE `playbot_fav` (
  `user` varchar(255) NOT NULL,
  `id` int(6) NOT NULL,
  `date` timestamp NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `playbot_later`
--

CREATE TABLE `playbot_later` (
  `id` int(11) NOT NULL,
  `content` int(11) NOT NULL,
  `nick` varchar(255) NOT NULL,
  `date` int(10) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `playbot_playlist_content_association`
--

CREATE TABLE `playbot_playlist_content_association` (
  `id` int(11) NOT NULL,
  `playlist_id` int(11) DEFAULT NULL,
  `content_id` int(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1 COLLATE=latin1_swedish_ci;

-- --------------------------------------------------------

--
-- Table structure for table `playbot_tags`
--

CREATE TABLE `playbot_tags` (
  `id` int(11) NOT NULL,
  `tag` varchar(50) CHARACTER SET latin1 COLLATE latin1_swedish_ci NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

--
-- Indexes for dumped tables
--

--
-- Indexes for table `playbot`
--
ALTER TABLE `playbot`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `url` (`url`),
  ADD KEY `chan` (`channel`);
ALTER TABLE `playbot` ADD FULLTEXT KEY `sender` (`sender`,`title`);

--
-- Indexes for table `playbot_chan`
--
ALTER TABLE `playbot_chan`
  ADD PRIMARY KEY (`id`),
  ADD KEY `content` (`content`),
  ADD KEY `chan` (`chan`);

--
-- Indexes for table `playbot_codes`
--
ALTER TABLE `playbot_codes`
  ADD PRIMARY KEY (`user`),
  ADD UNIQUE KEY `code` (`code`),
  ADD UNIQUE KEY `nick` (`nick`);

--
-- Indexes for table `playbot_config`
--
ALTER TABLE `playbot_config`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `playbot_fav`
--
ALTER TABLE `playbot_fav`
  ADD PRIMARY KEY (`user`,`id`),
  ADD KEY `id` (`id`),
  ADD KEY `user` (`user`);

--
-- Indexes for table `playbot_later`
--
ALTER TABLE `playbot_later`
  ADD PRIMARY KEY (`id`),
  ADD KEY `content` (`content`);

--
-- Indexes for table `playbot_playlist_content_association`
--
ALTER TABLE `playbot_playlist_content_association`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `playlist_id` (`playlist_id`,`content_id`);

--
-- Indexes for table `playbot_tags`
--
ALTER TABLE `playbot_tags`
  ADD PRIMARY KEY (`id`,`tag`),
  ADD KEY `tag` (`tag`),
  ADD KEY `id` (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `playbot`
--
ALTER TABLE `playbot`
  MODIFY `id` int(6) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `playbot_chan`
--
ALTER TABLE `playbot_chan`
  MODIFY `id` int(10) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `playbot_config`
--
ALTER TABLE `playbot_config`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `playbot_later`
--
ALTER TABLE `playbot_later`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `playbot_playlist_content_association`
--
ALTER TABLE `playbot_playlist_content_association`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `playbot_chan`
--
ALTER TABLE `playbot_chan`
  ADD CONSTRAINT `playbot_chan_ibfk_1` FOREIGN KEY (`content`) REFERENCES `playbot` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `playbot_fav`
--
ALTER TABLE `playbot_fav`
  ADD CONSTRAINT `playbot_fav_ibfk_2` FOREIGN KEY (`id`) REFERENCES `playbot` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `playbot_tags`
--
ALTER TABLE `playbot_tags`
  ADD CONSTRAINT `playbot_tags_ibfk_1` FOREIGN KEY (`id`) REFERENCES `playbot` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
