CREATE TABLE `undo_log`
(
    `pk_id`            bigint(20) AUTO_INCREMENT PRIMARY KEY,
    `xid`              varchar(100) NOT NULL,
    `branch_id`        bigint(20)   NOT NULL,
    `undo_items`       longblob     NOT NULL,
    `log_status`       int(11)      NOT NULL,
    `create_timestamp` bigint(20)   NOT NULL,
    UNIQUE KEY `ux_undo_log` (`xid`, `branch_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8;