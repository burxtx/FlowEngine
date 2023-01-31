## 轻量级工作流引擎
可以通过简单几行代码为系统增加审批功能。



## features

支持按条件自动审批

☑️  审批流节点无数量限制

☑️  摒弃复杂的UI，仅通过简单几行代码实现审批功能

☑️  支持内存流和持久化流

✅ 支持其他工作流（TODO）


  

## 使用方法

以创建持久化审批流为例来介绍：

0. 初始化审批流引擎

   ```
   import (
   	"github.com/burxtx/FlowEngine/approval"
   )
   
   f := approval.NewEngine(persistent.GetDB()) // 需传入 *gorm.DB 类型的数据库连接
   ```

   

1. 定义流程节点名

   需要显示指定开始和结束节点，默认第一个元素为开始节点，最后一个元素为结束节点

   ```
   nodes := []string{
   		"创建工单",
   		"主管审批",
   		"总监审批",
   		"总经理审批",
   		...  //还可以更多
   		"结束",
   	}
   processName := "上线提测流程" // 流程名称
   ```

   

2. 定义各审批人，需要是个二维数组

   ```
   approvers := [][]string{
   		{"Allen", "Penny"}, //一级审批人
   		{"Sussie"}, //二级审批人
   		{"Mikko"}, //三级审批人
   		... // 还可以更多
   }
   submitter := "chris" //申请人
   ```

   

3. 创建审批流

   ```
   fi, err := f.Create(ctx, approvers, nodes, submitter, processName)
   ```

   

4. 初始化调度器 

   ```
   r := scheduler.NewRule(price, GetQuota, amount, continuous, QuickApproval) // 初始化规则对象
   tx := db.GetDB(ctx)
   s := scheduler.NewApprovalScheduler(tx, r) // 初始化审批调度器
   ```

5. 执行审批

   ```
   err := s.Trigger(ctx, "approver_xxx", "pass", "memo", fi) 
   ```

   

   完整使用方法请参考单元测试用例

   

## 数据库表结构

用于审批流的持久化

### process 流程表

```
+------------+-------------+------+-----+---------+----------------+
| Field      | Type        | Null | Key | Default | Extra          |
+------------+-------------+------+-----+---------+----------------+
| id         | int(11)     | NO   | PRI | NULL    | auto_increment |
| service    | varchar(20) | NO   | MUL | NULL    |                |
| biz_id     | int(11)     | NO   |     | NULL    |                |
| head_node  | int(11)     | NO   |     | NULL    |                |
| tail_node  | int(11)     | NO   |     | NULL    |                |
| cur_state  | varchar(20) | NO   |     | NULL    |                |
| cur_node   | int(11)     | NO   |     | NULL    |                |
| closed_at  | int(11)     | NO   |     | NULL    |                |
| updated_at | int(11)     | NO   |     | NULL    |                |
| created_at | int(11)     | NO   | MUL | NULL    |                |
+------------+-------------+------+-----+---------+----------------+
```



### task_node 流程节点表

```
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| id          | int(11)      | NO   | PRI | NULL    | auto_increment |
| flow_id     | int(11)      | NO   |     | NULL    |                |
| role        | varchar(50)  | NO   | MUL | NULL    |                |
| value       | varchar(20)  | NO   | MUL | NULL    |                |
| remark      | varchar(500) | NO   |     | NULL    |                |
| status      | varchar(20)  | NO   | MUL | NULL    |                |
| result      | varchar(20)  | NO   |     | NULL    |                |
| pre_id      | int(11)      | NO   |     | NULL    |                |
| next_id     | int(11)      | NO   |     | NULL    |                |
| finished_at | int(11)      | NO   |     | NULL    |                |
| updated_at  | int(11)      | NO   |     | NULL    |                |
| created_at  | int(11)      | NO   | MUL | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
```



### node_line 节点流转表

```
+------------+---------+------+-----+---------+----------------+
| Field      | Type    | Null | Key | Default | Extra          |
+------------+---------+------+-----+---------+----------------+
| id         | int(11) | NO   | PRI | NULL    | auto_increment |
| flow_id    | int(11) | NO   |     | NULL    |                |
| parent     | int(11) | NO   |     | NULL    |                |
| child      | int(11) | NO   |     | NULL    |                |
| created_at | int(11) | NO   | MUL | NULL    |                |
| updated_at | int(11) | NO   |     | NULL    |                |
+------------+---------+------+-----+---------+----------------+
```



### candidate 审批人表

```
+------------+-------------+------+-----+---------+----------------+
| Field      | Type        | Null | Key | Default | Extra          |
+------------+-------------+------+-----+---------+----------------+
| id         | int(11)     | NO   | PRI | NULL    | auto_increment |
| node_id    | int(11)     | NO   | MUL | NULL    |                |
| approver   | varchar(20) | NO   | MUL | NULL    |                |
| updated_at | int(11)     | NO   |     | NULL    |                |
| created_at | int(11)     | NO   | MUL | NULL    |                |
| flow_id    | int(11)     | NO   |     | NULL    |                |
+------------+-------------+------+-----+---------+----------------+
```



### schedule_queue 审批队列表

```
+------------+--------------+------+-----+---------+----------------+
| Field      | Type         | Null | Key | Default | Extra          |
+------------+--------------+------+-----+---------+----------------+
| id         | int(11)      | NO   | PRI | NULL    | auto_increment |
| process_id | int(11)      | NO   | MUL | NULL    |                |
| node_id    | int(11)      | NO   |     | NULL    |                |
| user       | varchar(20)  | NO   |     | NULL    |                |
| name       | varchar(20)  | NO   |     | NULL    |                |
| memo       | varchar(20)  | NO   |     | NULL    |                |
| data       | varchar(500) | NO   |     | NULL    |                |
| created_at | int(11)      | NO   | MUL | NULL    |                |
| deleted_at | varchar(50)  | YES  |     | NULL    |                |
| state      | varchar(20)  | NO   |     | NULL    |                |
+------------+--------------+------+-----+---------+----------------+
```

