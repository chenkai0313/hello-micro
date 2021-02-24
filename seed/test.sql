
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
--  Table structure for `hello`
-- ----------------------------
DROP TABLE IF EXISTS `hello`;
CREATE TABLE `hello` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `content` text NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='微服务，demo 数据库';

-- ----------------------------
--  Records of `hello`
-- ----------------------------
BEGIN;
INSERT INTO `hello` VALUES ('25', 'hello', 'tom', '2021-02-23 06:25:05', '2021-02-23 06:25:05'), ('26', 'hello', 'tom', '2021-02-23 06:25:05', '2021-02-23 06:25:05');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
