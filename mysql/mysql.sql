-- 代理地址
DROP TABLE IF EXISTS `agents`;
CREATE TABLE IF NOT EXISTS `agents` (
  `id`          INT          NOT NULL AUTO_INCREMENT,
  `listen_ip`   VARCHAR(64)  NOT NULL,
  `listen_port` INT          NOT NULL,
  `desc`        VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
);

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
  `pause_flag`   INT          NOT NULL,
  `zk_address`   VARCHAR(256) NULL,
  `zk_node`      VARCHAR(256) NULL,
  `desc`         VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
);

-- 服务类型
DROP TABLE IF EXISTS `service_type`;
CREATE TABLE IF NOT EXISTS `service_type` (
  `id`           INT           NOT NULL AUTO_INCREMENT,
  `service_type` INT           NOT NULL UNIQUE,
  `service_name` VARCHAR(64)   NOT NULL,
  `start_cmd`    VARCHAR(1024) NOT NULL,
  `stop_cmd`     VARCHAR(1024) NOT NULL,
  `restart_cmd`  VARCHAR(1024) NULL,
  `desc`         VARCHAR(256)  NULL,
  PRIMARY KEY (`id`)
);

-- 服务依赖关系
DROP TABLE IF EXISTS `service_rely`;
CREATE TABLE IF NOT EXISTS `service_rely` (
  `id`      INT          NOT NULL AUTO_INCREMENT,
  `rely_id` INT          NOT NULL,
  `desc`    VARCHAR(256) NULL,
  PRIMARY KEY (`id`)
);

-- 用户账号
DROP TABLE IF EXISTS `users`;
CREATE TABLE IF NOT EXISTS `users` (
  `id`         INT          NOT NULL AUTO_INCREMENT,
  `email`      VARCHAR(256) NOT NULL,
  `login_type` INT          NOT NULL,
  `passwd`     VARCHAR(256) NULL,
  `right`      VARCHAR(64)  NOT NULL,
  `desc`       VARCHAR(256) NOT NULL,
  PRIMARY KEY (`id`)
);