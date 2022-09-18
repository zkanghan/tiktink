/*
SQLyog Community v13.1.6 (64 bit)
MySQL - 5.7.33 : Database - tiktink
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`tiktink` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_bin */;

USE `tiktink`;

/*Table structure for table `comments` */

DROP TABLE IF EXISTS `comments`;

CREATE TABLE `comments`
(
    `id`          bigint(20)                        NOT NULL AUTO_INCREMENT,
    `comment_id`  varchar(64) COLLATE utf8mb4_bin   NOT NULL COMMENT '评论编号',
    `author_id`   varchar(64) COLLATE utf8mb4_bin   NOT NULL COMMENT '作者编号',
    `video_id`    varchar(64) COLLATE utf8mb4_bin   NOT NULL COMMENT '视频编号',
    `content`     varchar(1000) COLLATE utf8mb4_bin NOT NULL COMMENT '评论内容',
    `create_date` timestamp                         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建日期',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_comment_id` (`comment_id`),
    KEY `idx_video_id` (`video_id`),
    KEY `idx_author_id` (`author_id`)
) ENGINE=InnoDB AUTO_INCREMENT=105 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*Table structure for table `favorites` */

DROP TABLE IF EXISTS `favorites`;

CREATE TABLE `favorites`
(
    `id`          bigint(20)                      NOT NULL AUTO_INCREMENT,
    `video_id`    varchar(64) COLLATE utf8mb4_bin NOT NULL,
    `user_id`     varchar(64) COLLATE utf8mb4_bin NOT NULL,
    `create_time` timestamp                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_video_id` (`video_id`),
    KEY `idx_author_id` (`user_id`, `video_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*Table structure for table `follow` */

DROP TABLE IF EXISTS `follow`;

CREATE TABLE `follow`
(
    `id`          bigint(20)                      NOT NULL AUTO_INCREMENT,
    `user_id`     varchar(64) COLLATE utf8mb4_bin NOT NULL,
    `to_user_id`  varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '被关注者的user_id',
    `create_time` timestamp                       NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`, `to_user_id`),
    KEY `idx_to_user_id` (`to_user_id`, `user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*Table structure for table `users` */

DROP TABLE IF EXISTS `users`;

CREATE TABLE `users`
(
    `id`             bigint(20)  NOT NULL AUTO_INCREMENT,
    `user_id`        varchar(64)          DEFAULT NULL,
    `user_name`      varchar(64) NOT NULL,
    `password`       varchar(64) NOT NULL,
    `follow_count`   bigint(20)  NOT NULL DEFAULT '0',
    `follower_count` bigint(20)  NOT NULL DEFAULT '0',
    `created_at`     timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at`     timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`     timestamp   NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_user_name` (`user_name`),
    UNIQUE KEY `uniq_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8mb4;

/*Table structure for table `videos` */

DROP TABLE IF EXISTS `videos`;

CREATE TABLE `videos`
(
    `id`             bigint(20)                       NOT NULL AUTO_INCREMENT,
    `author_id`      varchar(64) COLLATE utf8mb4_bin  NOT NULL COMMENT '作者编号',
    `play_url`       varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '视频在存储桶中的唯一标识',
    `video_id`       varchar(64) COLLATE utf8mb4_bin  NOT NULL COMMENT '视频编号',
    `cover_url`      varchar(200) COLLATE utf8mb4_bin NOT NULL COMMENT '封面路径',
    `image_id`       varchar(50) COLLATE utf8mb4_bin  NOT NULL COMMENT '封面编号',
    `favorite_count` bigint(20)                       NOT NULL DEFAULT '0' COMMENT '点赞数',
    `comment_count`  bigint(20)                       NOT NULL DEFAULT '0' COMMENT '评论数',
    `title`          varchar(20) COLLATE utf8mb4_bin  NOT NULL COMMENT '标题',
    `create_time`    timestamp                        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uniq_video_id` (`video_id`),
    KEY `idx_create_time` (`create_time`),
    KEY `idx_author_id` (`author_id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
