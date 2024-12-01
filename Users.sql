USE sso;
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
    `userID` varchar(13) NOT NULL,
    `password` varchar(50) NOT NULL,
    `name` varchar(20) NOT NULL,
    `grade` varchar(50) NULL,
    `majorClass` varchar(50) NULL,
    `email` varchar(50) NULL,
    `phoneNumber` varchar(13) NULL,
    `status` varchar(20) NULL,
    `role` varchar(20) NULL,
    PRIMARY KEY (`userID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `users` (`userID`, `password`, `name`) VALUES
  ('20233802086', 'qweasd','黄智泓'),
  ('20223808888', 'qweasd','万睿杰'),
  ('20241111111', 'qweasd','刘恒'),
  ('20231111111', 'qweasd','丁浩斌'),
  ('20223722222','qweasd','谢宇超');