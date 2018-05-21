-- 用户账号
DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id`          INT          NOT NULL AUTO_INCREMENT,
  `user_name`   VARCHAR(256) NULL,
  `user_email`  VARCHAR(256) NOT NULL,
  `user_pwd`    VARCHAR(256) NULL,
  `user_right`  BIGINT       NOT NULL DEFAULT 0,
  `user_type`   INT          NOT NULL DEFAULT 0,
  `create_time` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`        VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;
INSERT INTO users (user_name, user_email, user_pwd, user_right, user_type)
VALUES ('李艳国', 'lyg@meitu.com', 'e10adc3949ba59abbe56e057f20f883e', 1, 1);

-- 服务列表
DROP TABLE IF EXISTS `services`;
CREATE TABLE IF NOT EXISTS `services` (
  `id`           INT           NOT NULL AUTO_INCREMENT,
  `service_name` VARCHAR(256)  NOT NULL,
  `ssh_user`     VARCHAR(256)  NOT NULL,
  `ssh_ip`       VARCHAR(64)   NOT NULL,
  `start_cmd`    VARCHAR(1024) NOT NULL,
  `stop_cmd`     VARCHAR(1024) NOT NULL,
  `activate`     INT           NOT NULL DEFAULT 0,
  `auto_recover` INT           NOT NULL DEFAULT 1,
  `create_time`  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`         VARCHAR(256)  NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `service_name` (`service_name`)
) DEFAULT CHARSET = utf8;

-- 命令列表
DROP TABLE IF EXISTS `check_cmd`;
CREATE TABLE IF NOT EXISTS `check_cmd` (
  `id`              INT           NOT NULL AUTO_INCREMENT,
  `service_id`      INT           NOT NULL,
  `local_check`     INT           NOT NULL DEFAULT 0,
  `check_shell`     VARCHAR(1024) NOT NULL,
  `check_value`     VARCHAR(64)   NOT NULL,
  `health_if_match` INT           NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;

-- 依赖关系
DROP TABLE IF EXISTS `service_rely`;
CREATE TABLE IF NOT EXISTS `service_rely` (
  `id`          INT          NOT NULL AUTO_INCREMENT,
  `service_id`  INT          NULL,
  `create_time` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`        VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;