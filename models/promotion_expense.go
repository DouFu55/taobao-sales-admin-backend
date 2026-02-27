package models

import (
	"gorm.io/gorm"
	"time"
)

// PromotionExpense 推广支出记账模型（简化版）
type PromotionExpense struct {
	ID            int64     `gorm:"primaryKey;column:id" json:"id"`                                                    // 主键ID
	PromotionDate time.Time `gorm:"column:promotion_date;type:datetime;not null" json:"promotion_date"`                // 推广日期
	PromotionType int       `gorm:"column:promotion_type;type:tinyint;not null;default:1" json:"promotion_type"`       // 推广类型：1=全站推广，2=关键词推广，3=人群推广，4=内容推广，5=淘宝联盟，6=营销托管，7=其他
	PlanID        string    `gorm:"column:plan_id;type:varchar(100)" json:"plan_id"`                                   // 计划ID/名称
	ExpenseAmount float64   `gorm:"column:expense_amount;type:decimal(10,2);not null;default:0" json:"expense_amount"` // 支出金额(元)
	OrdersCount   int       `gorm:"column:orders_count;default:0" json:"orders_count"`                                 // 成交笔数
	Remark        string    `gorm:"column:remark;type:varchar(200)" json:"remark"`                                     // 备注

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (PromotionExpense) TableName() string {
	return "promotion_expenses"
}
