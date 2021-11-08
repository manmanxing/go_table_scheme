# go_table_scheme
SQL语句到结构体的转换工具

# 使用方法
在本项目根目录下或 go install 后执行下面命令
go_table_scheme -t 表名 -db '用户名:密码@tcp(ip:3306)/数据库名'

# 示例
数据库信息
```mysql
CREATE TABLE `trade` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
    `user_id` bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT '用户id',
    `trade_no` varchar(64) NOT NULL DEFAULT '' COMMENT '交易号',
    `pay_no` varchar(64) NOT NULL DEFAULT '' COMMENT '支付单号',
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `trade_pay_no_uindex` (`pay_no`),
    UNIQUE KEY `trade_trade_no_uindex` (`trade_no`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 ROW_FORMAT = DYNAMIC COMMENT = '交易主表'
```

转换成结构体如下
```go
package model
import (
	_ "fmt"
	"time"
)

type Trade struct {
	ID        int64     `sql:"primary_key;column:id" json:"id,omitempty"`     //主键
	UserID    int64     `sql:"column:user_id" json:"user_id,omitempty"`       //用户id
	TradeNO   string    `sql:"column:trade_no" json:"trade_no,omitempty"`     //交易号
	PayNO     string    `sql:"column:pay_no" json:"pay_no,omitempty"`         //支付单号
	CreatedAt time.Time `sql:"column:created_at" json:"created_at,omitempty"` //创建时间
	UpdatedAt time.Time `sql:"column:updated_at" json:"updated_at,omitempty"` //更新时间

}

func (Trade) TableName() string {
	return "trade"
}
```