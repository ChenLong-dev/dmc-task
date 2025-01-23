-- goctl model mysql ddl --src t_cron_cycle_tasks.sql --dir .

CREATE TABLE IF NOT EXISTS `t_cron_cycle_tasks`
(
    `id`          char(128) UNIQUE NOT NULL COMMENT '定时循环任务ID',
    `entry_id`    tinyint(32)      NOT NULL DEFAULT 0 COMMENT '入口ID',
    `type`        tinyint(32)      NOT NULL DEFAULT 0 COMMENT '任务类型',
    `biz_code`    varchar(128)     NOT NULL DEFAULT '' COMMENT '业务Code',
    `cron`        varchar(128)     NOT NULL DEFAULT '' COMMENT 'cron参数',
    `exec_path`   varchar(1024)    NOT NULL DEFAULT '' COMMENT '执行路径',
    `param`       varchar(1024)    NOT NULL DEFAULT '' COMMENT '任务的执行参数',
    `timeout`     int(11)          NOT NULL DEFAULT 0 COMMENT '任务超时时间，单位秒',
    `status`      tinyint          NOT NULL DEFAULT 0 COMMENT '任务执行状态',
    `ext_info`    json             NOT NULL COMMENT '扩展信息',
    `update_time` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '任务更新时间',
    `create_time` datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '任务创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = Dynamic COMMENT ='任务配置表';