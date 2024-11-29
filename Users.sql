CREATE DATABASE sso;
USE sso;
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `userID` varchar(13) NOT NULL,
    `name` varchar(20) NOT NULL,
    PRIMARY KEY (`userID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `users` VALUES
  ('20233802086', '黄智泓'),
  ('20223808888', '万睿杰'),
  ('20241111111', '刘恒'),
  ('20231111111', '丁浩斌');