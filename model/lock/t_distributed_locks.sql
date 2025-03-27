-- goctl model mysql ddl --src t_distributed_locks.sql --dir .

CREATE TABLE IF NOT EXISTS t_distributed_locks
(
    `lock_name`   VARCHAR(64)  NOT NULL,
    `source`      VARCHAR(128) NOT NULL,
    `lock_value`  VARCHAR(128) NOT NULL,
    `expire_time` TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`lock_name`),
    UNIQUE KEY `idx_lock_name` (`lock_name`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = Dynamic COMMENT ='分布式锁表';
