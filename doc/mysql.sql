-- 代理地址
DROP TABLE IF EXISTS `agents`;
CREATE TABLE IF NOT EXISTS `agents` (
  `id`          INT          NOT NULL AUTO_INCREMENT,
  `listen_ip`   VARCHAR(64)  NOT NULL,
  `listen_port` INT          NOT NULL,
  `create_time` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`        VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;

-- 服务安装信息
DROP TABLE IF EXISTS `services`;
CREATE TABLE IF NOT EXISTS `services` (
  `id`           INT          NOT NULL AUTO_INCREMENT,
  `service_type` INT          NOT NULL,
  `ip_address`   VARCHAR(64)  NOT NULL,
  `domain_name`  VARCHAR(256) NULL,
  `install_path` VARCHAR(256) NOT NULL,
  `log_path`     VARCHAR(256) NOT NULL,
  `run_user`     VARCHAR(64)  NOT NULL,
  `listen_port`  INT          NULL,
  `pause_flag`   INT          NOT NULL DEFAULT 0,
  `zk_address`   VARCHAR(256) NULL,
  `zk_node`      VARCHAR(256) NULL,
  `create_time`  TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`         VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;

-- 服务类型
DROP TABLE IF EXISTS `service_type`;
CREATE TABLE IF NOT EXISTS `service_type` (
  `id`           INT           NOT NULL AUTO_INCREMENT,
  `service_type` INT           NOT NULL UNIQUE,
  `service_name` VARCHAR(64)   NOT NULL,
  `start_cmd`    VARCHAR(1024) NOT NULL,
  `stop_cmd`     VARCHAR(1024) NOT NULL,
  `restart_cmd`  VARCHAR(1024) NULL,
  `create_time`  TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`         VARCHAR(256)  NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;

-- 服务依赖关系
DROP TABLE IF EXISTS `service_rely`;
CREATE TABLE IF NOT EXISTS `service_rely` (
  `id`          INT          NOT NULL AUTO_INCREMENT,
  `rely_id`     INT          NULL,
  `create_time` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `desc`        VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
) DEFAULT CHARSET = utf8;

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