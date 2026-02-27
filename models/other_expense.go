package models

import (
	"gorm.io/gorm"
	"time"
)

// OtherExpense 其他支出记账模型（简化版）
type OtherExpense struct {
	ID          int64     `gorm:"primaryKey;column:id" json:"id"`                                          // 主键ID
	ExpenseDate time.Time `gorm:"column:expense_date;type:datetime;not null" json:"expense_date"`          // 支出日期
	ExpenseType int       `gorm:"column:expense_type;type:tinyint;not null;default:1" json:"expense_type"` // 支出类型：1=办公用品，2=水电费，3=工资，4=交通费，5=差旅费，6=设备购置，7=其他
	ItemName    string    `gorm:"column:item_name;type:varchar(200);not null" json:"item_name"`            // 项目名称
	Amount      float64   `gorm:"column:amount;type:decimal(10,2);not null;default:0" json:"amount"`       // 金额(元)
	SubOrderNo  string    `gorm:"column:sub_order_no;type:varchar(100);index" json:"sub_order_no"`         // 子订单号（选填）
	Remark      string    `gorm:"column:remark;type:varchar(500)" json:"remark"`                           // 备注

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`          // 删除时间
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"` // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` // 更新时间
}

func (OtherExpense) TableName() string {
	return "other_expenses"
}
