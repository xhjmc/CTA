CREATE TABLE `test`
(
    `pk_id` bigint(20) AUTO_INCREMENT PRIMARY KEY,
    `id`    int(11)      NOT NULL,
    `col`   varchar(100) NOT NULL
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4;

truncate test;