
CREATE TABLE `process` (
`id`          int(11)      NOT NULL AUTO_INCREMENT COMMENT "自增id",
`service`     varchar(20)  NOT NULL                COMMENT "流程名",
`biz_id` int(11) NOT NULL COMMENT "业务ID",
`head_node` int(11) NOT NULL COMMENT "头节点ID",
`tail_node` int(11) NOT NULL COMMENT "尾节点ID",
`cur_state` varchar(20) NOT NULL,
`cur_node`      int(11)      NOT NULL                COMMENT "当前节点id",
`closed_at`   int(11)      NOT NULL                COMMENT "关闭工单时间",
`updated_at`     int(11)      NOT NULL                COMMENT "修改时间",
`created_at`  int(11)      NOT NULL                COMMENT "创建工单时间",
PRIMARY KEY           (`id`),
KEY `idx_service`     (`service`),
KEY `idx_created`     (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;


CREATE TABLE `task_node` (
`id`             int(11)      NOT NULL AUTO_INCREMENT COMMENT "自增id",
`flow_id`      int(11)      NOT NULL                COMMENT "工单id",
`role`           varchar(50)  NOT NULL                COMMENT "审批角色",
`value`       varchar(20)  NOT NULL                COMMENT "最终审批人ID",
`remark`         varchar(500) NOT NULL                COMMENT "审批评语",
`status`         varchar(20)   NOT NULL                COMMENT "当前审批状态",
`result` varchar(20) NOT NULL COMMENT "结果",
`pre_id` int(11) NOT NULL,
`next_id` int(11) NOT NULL,
`finished_at`    int(11)      NOT NULL                COMMENT "完成工单时间",
`updated_at`     int(11)      NOT NULL                COMMENT "修改时间",
`created_at`     int(11)      NOT NULL                COMMENT "创建工单时间",
PRIMARY KEY                     (`id`),
KEY        `idx_role`           (`role`),
KEY        `idx_value`       (`value`),
KEY        `idx_status`         (`status`),
KEY        `idx_created_at`     (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;


CREATE TABLE `node_line` (
`id`             int(11)     NOT NULL AUTO_INCREMENT COMMENT "自增id",
`flow_id`      int(11)     NOT NULL                COMMENT "工单id",
`parent`       int(11) NOT NULL                COMMENT "父节点",
`child`       int(11) NOT NULL                COMMENT "子节点",
`created_at`     int(11)     NOT NULL                COMMENT "创建时间",
`updated_at`     int(11)      NOT NULL                COMMENT "修改时间",
PRIMARY KEY                   (`id`),
KEY `idx_created_at`          (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE `candidate` (
`id`             int(11)     NOT NULL AUTO_INCREMENT COMMENT "自增id",
`node_id` int(11)     NOT NULL                COMMENT "节点id",
`flow_id` int(11)     NOT NULL                COMMENT "工单id",
`approver`       varchar(20) NOT NULL                COMMENT "审批人ID",
`updated_at`     int(11)     NOT NULL                COMMENT "修改时间",
`created_at`     int(11)     NOT NULL                COMMENT "创建时间",
PRIMARY KEY                   (`id`),
KEY `idx_node_id`      (`node_id`),
KEY `idx_approver`            (`approver`),
KEY `idx_created_at`          (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;


CREATE TABLE `change_log` (
`id`             int(11)      NOT NULL AUTO_INCREMENT COMMENT "自增id",
`process_id`      int(11)      NOT NULL                COMMENT "工单id",
`user`       varchar(20)  NOT NULL                COMMENT "执行者",
`auto_skip`   tinyint(8)   NOT NULL                COMMENT "是否自动跳过",
`data`         varchar(500) NOT NULL                COMMENT "审批数据",
`created_at`     int(11)      NOT NULL                COMMENT "操作时间",
PRIMARY KEY           (`id`),
KEY `idx_process_id`   (`process_id`),
KEY `idx_created_at`  (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;


CREATE TABLE `schedule_queue` (
`id`             int(11)      NOT NULL AUTO_INCREMENT COMMENT "自增id",
`process_id`      int(11)      NOT NULL                COMMENT "工单id",
`node_id`      int(11)      NOT NULL                COMMENT "节点id",
`user`       varchar(20)  NOT NULL                COMMENT "执行者",
`state`       varchar(20)  NOT NULL                COMMENT "审批状态",
`name`       varchar(20)  NOT NULL                COMMENT "节点名",
`memo`       varchar(20)  NOT NULL                COMMENT "备注",
`data`         varchar(500) NOT NULL                COMMENT "审批数据",
`created_at`     int(11)      NOT NULL                COMMENT "操作时间",
`deleted_at`     varchar(50)      NULL                COMMENT "删除时间",
PRIMARY KEY           (`id`),
KEY `idx_process_id`   (`process_id`),
KEY `idx_created_at`  (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

