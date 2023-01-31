use DATABASE;
alter table flow
modify `service` varchar(20) NOT NULL COMMENT "流程名",
modify `description` varchar(500) NOT NULL DEFAULT '';
add column `biz_id` int(11) NOT NULL COMMENT "业务ID",
add column `head_node` int(11) NOT NULL COMMENT "头节点ID",
add column `tail_node` int(11) NOT NULL COMMENT "尾节点ID",
add column `cur_state` varchar(20) NOT NULL COMMENT "当前状态",
add column `cur_state_cn_name` varchar(20) NOT NULL COMMENT "当前状态中文",
add column `updated_at` int(11) NOT NULL COMMENT "修改时间";

alter table node
modify column `flow_id` int(11) NOT NULL COMMENT "工单id",
modify column `role` varchar(20) NOT NULL,
add column `pre_id` int(11) NOT NULL COMMENT "前一节点ID",
add column `next_id` int(11) NOT NULL COMMENT "后一节点ID",
add column `result` varchar(20) NOT NULL COMMENT "结果",
change column `update_time` `updated_at` int(11) NOT NULL,
change column `finish_time` `finished_at` int(11) NOT NULL

alter table node_relation
add column `updated_at` int(11) NOT NULL COMMENT "修改时间";

CREATE TABLE candidate (
`id`             int(11)     NOT NULL AUTO_INCREMENT COMMENT "自增id",
`node_id` int(11)     NOT NULL                COMMENT "节点id",
`flow_id` int(11)     NOT NULL                COMMENT "工单id",
`approver`       varchar(20) NOT NULL                COMMENT "审批人邮箱前缀",
`updated_at`     int(11)     NOT NULL                COMMENT "修改时间",
`created_at`     int(11)     NOT NULL                COMMENT "创建时间",
PRIMARY KEY                   (`id`),
KEY `idx_node_id`      (`node_id`),
KEY `idx_approver`            (`approver`),
KEY `idx_created_at`          (`created_at`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;
