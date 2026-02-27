package models

import (
	"gorm.io/gorm"
	"time"
)

// SalesItemDetail 订单模型（使用指针类型避免空值问题）
type SalesItemDetail struct {
	gorm.Model
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`                                             // 主键ID
	SubOrderID      *string        `gorm:"column:sub_order_id;type:varchar(50);uniqueIndex" json:"sub_order_id,omitempty"` // 子订单号
	MainOrderID     *string        `gorm:"column:main_order_id;type:varchar(50);index" json:"main_order_id,omitempty"`     // 主订单号
	ProductTitle    *string        `gorm:"column:product_title;type:text" json:"product_title,omitempty"`                  // 商品标题
	ProductPrice    *float64       `gorm:"column:product_price;type:decimal(10,2)" json:"product_price,omitempty"`         // 商品单价
	Quantity        *int           `gorm:"column:quantity" json:"quantity,omitempty"`                                      // 购买数量
	ProductAttr     *string        `gorm:"column:product_attr;type:text" json:"product_attr,omitempty"`                    // 商品属性
	MealInfo        *string        `gorm:"column:meal_info;type:text" json:"meal_info,omitempty"`                          // 套餐信息
	OrderStatus     *string        `gorm:"column:order_status;type:varchar(100)" json:"order_status,omitempty"`            // 订单状态
	SkuCode         *string        `gorm:"column:sku_code;type:varchar(50);index" json:"sku_code,omitempty"`               // SKU编码
	BuyerPaid       *float64       `gorm:"column:buyer_paid;type:decimal(10,2)" json:"buyer_paid,omitempty"`               // 买家实付金额
	BuyerActualPay  *float64       `gorm:"column:buyer_actual_pay;type:decimal(10,2)" json:"buyer_actual_pay,omitempty"`   // 买家应付货款
	RefundStatus    *string        `gorm:"column:refund_status;type:varchar(50)" json:"refund_status,omitempty"`           // 退款状态
	RefundAmount    *float64       `gorm:"column:refund_amount;type:decimal(10,2)" json:"refund_amount,omitempty"`         // 退款金额
	OrderCreateTime *time.Time     `gorm:"column:order_create_time;index" json:"order_create_time,omitempty"`              // 订单创建时间
	OrderPayTime    *time.Time     `gorm:"column:order_pay_time;index" json:"order_pay_time,omitempty"`                    // 订单支付时间
	ProductID       *string        `gorm:"column:product_id;type:varchar(50)" json:"product_id,omitempty"`                 // 商品ID
	SellerRemark    *string        `gorm:"column:seller_remark;type:text" json:"seller_remark,omitempty"`                  // 卖家备注
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`                                      // 软删除时间
	CreatedAt       time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`                             // 创建时间
	UpdatedAt       time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`                             // 更新时间
}

// TableName 指定表名
func (SalesItemDetail) TableName() string {
	return "sales_item_detail"
}

// GetSubOrderID 安全获取子订单ID
func (o *SalesItemDetail) GetSubOrderID() string {
	if o.SubOrderID == nil {
		return ""
	}
	return *o.SubOrderID
}

// GetSkuCode 安全获取SKU编码
func (o *SalesItemDetail) GetSkuCode() string {
	if o.SkuCode == nil {
		return ""
	}
	return *o.SkuCode
}

// GetBuyerActualPay 安全获取买家应付货款
func (o *SalesItemDetail) GetBuyerActualPay() float64 {
	if o.BuyerActualPay == nil {
		return 0
	}
	return *o.BuyerActualPay
}
