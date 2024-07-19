-- 创建一个单独的用户
CREATE USER IF NOT EXISTS 'lzb200244'@'%' IDENTIFIED BY 'lzb200244';
GRANT ALL PRIVILEGES ON *.* TO 'lzb200244'@'%';
-- 创建数据库
CREATE DATABASE IF NOT EXISTS select_course CHARACTER SET utf8mb4;
/*
 Navicat Premium Dump SQL

 Source Server         : master
 Source Server Type    : MySQL
 Source Server Version : 80038 (8.0.38)
 Source Host           : 192.168.241.128:3306
 Source Schema         : select_course

 Target Server Type    : MySQL
 Target Server Version : 80038 (8.0.38)
 File Encoding         : 65001

 Date: 19/07/2024 22:10:55
*/

USE select_course;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for course
-- ----------------------------
DROP TABLE IF EXISTS `course`;
CREATE TABLE `course`
(
    `id`          bigint UNSIGNED                                              NOT NULL AUTO_INCREMENT,
    `title`       varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '课程名称',
    `category_id` bigint UNSIGNED                                              NOT NULL COMMENT '分类ID',
    `schedule_id` bigint UNSIGNED                                              NOT NULL COMMENT '分类ID',
    `capacity`    bigint                                                       NOT NULL COMMENT '容纳人数',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `fk_course_category` (`category_id` ASC) USING BTREE,
    INDEX `fk_course_schedule` (`schedule_id` ASC) USING BTREE,
    CONSTRAINT `fk_course_category` FOREIGN KEY (`category_id`) REFERENCES `course_category` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_course_schedule` FOREIGN KEY (`schedule_id`) REFERENCES `schedule` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  AUTO_INCREMENT = 11
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of course
-- ----------------------------
INSERT INTO `course`
VALUES (1, '课程1', 1, 8, 10);
INSERT INTO `course`
VALUES (2, '课程2', 3, 12, 10);
INSERT INTO `course`
VALUES (3, '课程3', 5, 13, 10);
INSERT INTO `course`
VALUES (4, '课程4', 3, 5, 10);
INSERT INTO `course`
VALUES (5, '课程5', 5, 9, 10);
INSERT INTO `course`
VALUES (6, '课程6', 4, 1, 10);
INSERT INTO `course`
VALUES (7, '课程7', 4, 8, 10);
INSERT INTO `course`
VALUES (8, '课程8', 4, 4, 10);
INSERT INTO `course`
VALUES (9, '课程9', 2, 1, 10);
INSERT INTO `course`
VALUES (10, '课程10', 3, 15, 10);

-- ----------------------------
-- Table structure for course_category
-- ----------------------------
DROP TABLE IF EXISTS `course_category`;
CREATE TABLE `course_category`
(
    `id`   bigint UNSIGNED                                              NOT NULL AUTO_INCREMENT,
    `name` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '分类名称',
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 6
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of course_category
-- ----------------------------
INSERT INTO `course_category`
VALUES (1, '分类1');
INSERT INTO `course_category`
VALUES (2, '分类2');
INSERT INTO `course_category`
VALUES (3, '分类3');
INSERT INTO `course_category`
VALUES (4, '分类4');
INSERT INTO `course_category`
VALUES (5, '分类5');

-- ----------------------------
-- Table structure for schedule
-- ----------------------------
DROP TABLE IF EXISTS `schedule`;
CREATE TABLE `schedule`
(
    `id`       bigint UNSIGNED NOT NULL AUTO_INCREMENT,
    `duration` tinyint         NULL DEFAULT NULL,
    `week`     tinyint         NULL DEFAULT NULL,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 16
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of schedule
-- ----------------------------
INSERT INTO `schedule`
VALUES (1, 1, 0);
INSERT INTO `schedule`
VALUES (2, 2, 0);
INSERT INTO `schedule`
VALUES (3, 3, 0);
INSERT INTO `schedule`
VALUES (4, 1, 1);
INSERT INTO `schedule`
VALUES (5, 2, 1);
INSERT INTO `schedule`
VALUES (6, 3, 1);
INSERT INTO `schedule`
VALUES (7, 1, 2);
INSERT INTO `schedule`
VALUES (8, 2, 2);
INSERT INTO `schedule`
VALUES (9, 3, 2);
INSERT INTO `schedule`
VALUES (10, 1, 3);
INSERT INTO `schedule`
VALUES (11, 2, 3);
INSERT INTO `schedule`
VALUES (12, 3, 3);
INSERT INTO `schedule`
VALUES (13, 1, 4);
INSERT INTO `schedule`
VALUES (14, 2, 4);
INSERT INTO `schedule`
VALUES (15, 3, 4);

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`
(
    `id`        bigint UNSIGNED                                              NOT NULL AUTO_INCREMENT,
    `user_name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '用户名称',
    `password`  varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
    `flag`      int                                                          NOT NULL COMMENT '用户标准位记录着选课的已选字段',
    PRIMARY KEY (`id`) USING BTREE,
    INDEX `idx_user_user_name` (`user_name` ASC) USING BTREE
) ENGINE = InnoDB
  AUTO_INCREMENT = 101
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user
-- ----------------------------
INSERT INTO `user`
VALUES (1, 'users1', 'password1', 0);
INSERT INTO `user`
VALUES (2, 'users2', 'password2', 0);
INSERT INTO `user`
VALUES (3, 'users3', 'password3', 0);
INSERT INTO `user`
VALUES (4, 'users4', 'password4', 0);
INSERT INTO `user`
VALUES (5, 'users5', 'password5', 0);
INSERT INTO `user`
VALUES (6, 'users6', 'password6', 0);
INSERT INTO `user`
VALUES (7, 'users7', 'password7', 0);
INSERT INTO `user`
VALUES (8, 'users8', 'password8', 0);
INSERT INTO `user`
VALUES (9, 'users9', 'password9', 0);
INSERT INTO `user`
VALUES (10, 'users10', 'password10', 0);
INSERT INTO `user`
VALUES (11, 'users11', 'password11', 0);
INSERT INTO `user`
VALUES (12, 'users12', 'password12', 0);
INSERT INTO `user`
VALUES (13, 'users13', 'password13', 0);
INSERT INTO `user`
VALUES (14, 'users14', 'password14', 0);
INSERT INTO `user`
VALUES (15, 'users15', 'password15', 0);
INSERT INTO `user`
VALUES (16, 'users16', 'password16', 0);
INSERT INTO `user`
VALUES (17, 'users17', 'password17', 0);
INSERT INTO `user`
VALUES (18, 'users18', 'password18', 0);
INSERT INTO `user`
VALUES (19, 'users19', 'password19', 0);
INSERT INTO `user`
VALUES (20, 'users20', 'password20', 0);
INSERT INTO `user`
VALUES (21, 'users21', 'password21', 0);
INSERT INTO `user`
VALUES (22, 'users22', 'password22', 0);
INSERT INTO `user`
VALUES (23, 'users23', 'password23', 0);
INSERT INTO `user`
VALUES (24, 'users24', 'password24', 0);
INSERT INTO `user`
VALUES (25, 'users25', 'password25', 0);
INSERT INTO `user`
VALUES (26, 'users26', 'password26', 0);
INSERT INTO `user`
VALUES (27, 'users27', 'password27', 0);
INSERT INTO `user`
VALUES (28, 'users28', 'password28', 0);
INSERT INTO `user`
VALUES (29, 'users29', 'password29', 0);
INSERT INTO `user`
VALUES (30, 'users30', 'password30', 0);
INSERT INTO `user`
VALUES (31, 'users31', 'password31', 0);
INSERT INTO `user`
VALUES (32, 'users32', 'password32', 0);
INSERT INTO `user`
VALUES (33, 'users33', 'password33', 0);
INSERT INTO `user`
VALUES (34, 'users34', 'password34', 0);
INSERT INTO `user`
VALUES (35, 'users35', 'password35', 0);
INSERT INTO `user`
VALUES (36, 'users36', 'password36', 0);
INSERT INTO `user`
VALUES (37, 'users37', 'password37', 0);
INSERT INTO `user`
VALUES (38, 'users38', 'password38', 0);
INSERT INTO `user`
VALUES (39, 'users39', 'password39', 0);
INSERT INTO `user`
VALUES (40, 'users40', 'password40', 0);
INSERT INTO `user`
VALUES (41, 'users41', 'password41', 0);
INSERT INTO `user`
VALUES (42, 'users42', 'password42', 0);
INSERT INTO `user`
VALUES (43, 'users43', 'password43', 0);
INSERT INTO `user`
VALUES (44, 'users44', 'password44', 0);
INSERT INTO `user`
VALUES (45, 'users45', 'password45', 0);
INSERT INTO `user`
VALUES (46, 'users46', 'password46', 0);
INSERT INTO `user`
VALUES (47, 'users47', 'password47', 0);
INSERT INTO `user`
VALUES (48, 'users48', 'password48', 0);
INSERT INTO `user`
VALUES (49, 'users49', 'password49', 0);
INSERT INTO `user`
VALUES (50, 'users50', 'password50', 0);
INSERT INTO `user`
VALUES (51, 'users51', 'password51', 0);
INSERT INTO `user`
VALUES (52, 'users52', 'password52', 0);
INSERT INTO `user`
VALUES (53, 'users53', 'password53', 0);
INSERT INTO `user`
VALUES (54, 'users54', 'password54', 0);
INSERT INTO `user`
VALUES (55, 'users55', 'password55', 0);
INSERT INTO `user`
VALUES (56, 'users56', 'password56', 0);
INSERT INTO `user`
VALUES (57, 'users57', 'password57', 0);
INSERT INTO `user`
VALUES (58, 'users58', 'password58', 0);
INSERT INTO `user`
VALUES (59, 'users59', 'password59', 0);
INSERT INTO `user`
VALUES (60, 'users60', 'password60', 0);
INSERT INTO `user`
VALUES (61, 'users61', 'password61', 0);
INSERT INTO `user`
VALUES (62, 'users62', 'password62', 0);
INSERT INTO `user`
VALUES (63, 'users63', 'password63', 0);
INSERT INTO `user`
VALUES (64, 'users64', 'password64', 0);
INSERT INTO `user`
VALUES (65, 'users65', 'password65', 0);
INSERT INTO `user`
VALUES (66, 'users66', 'password66', 0);
INSERT INTO `user`
VALUES (67, 'users67', 'password67', 0);
INSERT INTO `user`
VALUES (68, 'users68', 'password68', 0);
INSERT INTO `user`
VALUES (69, 'users69', 'password69', 0);
INSERT INTO `user`
VALUES (70, 'users70', 'password70', 0);
INSERT INTO `user`
VALUES (71, 'users71', 'password71', 0);
INSERT INTO `user`
VALUES (72, 'users72', 'password72', 0);
INSERT INTO `user`
VALUES (73, 'users73', 'password73', 0);
INSERT INTO `user`
VALUES (74, 'users74', 'password74', 0);
INSERT INTO `user`
VALUES (75, 'users75', 'password75', 0);
INSERT INTO `user`
VALUES (76, 'users76', 'password76', 0);
INSERT INTO `user`
VALUES (77, 'users77', 'password77', 0);
INSERT INTO `user`
VALUES (78, 'users78', 'password78', 0);
INSERT INTO `user`
VALUES (79, 'users79', 'password79', 0);
INSERT INTO `user`
VALUES (80, 'users80', 'password80', 0);
INSERT INTO `user`
VALUES (81, 'users81', 'password81', 0);
INSERT INTO `user`
VALUES (82, 'users82', 'password82', 0);
INSERT INTO `user`
VALUES (83, 'users83', 'password83', 0);
INSERT INTO `user`
VALUES (84, 'users84', 'password84', 0);
INSERT INTO `user`
VALUES (85, 'users85', 'password85', 0);
INSERT INTO `user`
VALUES (86, 'users86', 'password86', 0);
INSERT INTO `user`
VALUES (87, 'users87', 'password87', 0);
INSERT INTO `user`
VALUES (88, 'users88', 'password88', 0);
INSERT INTO `user`
VALUES (89, 'users89', 'password89', 0);
INSERT INTO `user`
VALUES (90, 'users90', 'password90', 0);
INSERT INTO `user`
VALUES (91, 'users91', 'password91', 0);
INSERT INTO `user`
VALUES (92, 'users92', 'password92', 0);
INSERT INTO `user`
VALUES (93, 'users93', 'password93', 0);
INSERT INTO `user`
VALUES (94, 'users94', 'password94', 0);
INSERT INTO `user`
VALUES (95, 'users95', 'password95', 0);
INSERT INTO `user`
VALUES (96, 'users96', 'password96', 0);
INSERT INTO `user`
VALUES (97, 'users97', 'password97', 0);
INSERT INTO `user`
VALUES (98, 'users98', 'password98', 0);
INSERT INTO `user`
VALUES (99, 'users99', 'password99', 0);
INSERT INTO `user`
VALUES (100, 'users100', 'password100', 0);

-- ----------------------------
-- Table structure for user_course
-- ----------------------------
DROP TABLE IF EXISTS `user_course`;
CREATE TABLE `user_course`
(
    `user_id`    bigint UNSIGNED NOT NULL COMMENT '用户ID',
    `course_id`  bigint UNSIGNED NOT NULL COMMENT '课程ID',
    `created_at` bigint          NULL     DEFAULT NULL,
    `updated_at` bigint          NULL     DEFAULT NULL,
    `is_deleted` tinyint(1)      NOT NULL DEFAULT 0 COMMENT '是否删除',
    UNIQUE INDEX `user_course` (`user_id` ASC, `course_id` ASC) USING BTREE,
    INDEX `fk_user_course_course` (`course_id` ASC) USING BTREE,
    CONSTRAINT `fk_user_course_course` FOREIGN KEY (`course_id`) REFERENCES `course` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT `fk_user_course_user` FOREIGN KEY (`user_id`) REFERENCES `user` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci
  ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_course
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;

FLUSH PRIVILEGES;