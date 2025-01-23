-- goctl model mysql ddl --src t_jobs_flow.sql --dir .

CREATE TABLE IF NOT EXISTS `t_jobs_flow`
(
    `id`            char(128) UNIQUE NOT NULL COMMENT 'job的任务ID',
    `type`          tinyint          NOT NULL DEFAULT 0 COMMENT '任务类型',
    `cron_task_id`  char(128)        NULL     DEFAULT '' COMMENT '定时任务ID',
    `biz_code`      varchar(128)     NOT NULL DEFAULT '' COMMENT '业务Code',
    `biz_id`        varchar(128)     NULL     DEFAULT '' COMMENT '业务ID',
    `exec_path`     varchar(1024)    NOT NULL DEFAULT '' COMMENT '执行路径',
    `param`         varchar(1024)    NOT NULL DEFAULT '' COMMENT '任务的执行参数',
    `timeout`       int(11)              NOT NULL DEFAULT 0 COMMENT '任务超时时间，单位秒',
    `start_time`    datetime         NULL     DEFAULT NULL COMMENT '定时任务执行的实际开始时间',
    `finish_time`   datetime         NULL     DEFAULT NULL COMMENT '定时任务执行的实际结束时间',
    `exec_interval` mediumint        NOT NULL DEFAULT 0 COMMENT 'job的执行时间（finish_time-start_time）',
    `status`        tinyint          NOT NULL DEFAULT 0 COMMENT '执行状态',
    `result_msg`    json             NOT NULL COMMENT 'Job执行结果的描述',
    `ext_info`      json             NOT NULL COMMENT '扩展信息',
    `update_time`   datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'job flow的更新时间',
    `create_time`   datetime         NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'job flow的创建时间',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_unicode_ci
  ROW_FORMAT = Dynamic COMMENT ='job流水记录表';