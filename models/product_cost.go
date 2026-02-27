package models

import (
	"gorm.io/gorm"
	"time"
)

// ProductCost 产品成本表
type ProductCost struct {
	ID uint `gorm:"primarykey" json:"id"`

	// SKU字段（必填，非空即可）
	SKU string `gorm:"type:varchar(100);not null;comment:商家编码" json:"sku"`

	// 商品信息
	ProductID   uint64 `gorm:"index;not null;comment:商品ID" json:"product_id"`
	ProductName string `gorm:"type:varchar(200);comment:商品名称" json:"product_name"`

	// 规格信息
	SpecInfo string `gorm:"type:varchar(500);comment:规格属性" json:"spec_info"`

	// 成本信息
	PurchasePrice float64 `gorm:"type:decimal(10,2);default:0.00;comment:采购价格" json:"purchase_price"`
	ShippingCost  float64 `gorm:"type:decimal(10,2);default:0.00;comment:运费" json:"shipping_cost"`
	HandlingCost  float64 `gorm:"type:decimal(10,2);default:0.00;comment:操作费" json:"handling_cost"`
	OtherCost     float64 `gorm:"type:decimal(10,2);default:0.00;comment:其他费用" json:"other_cost"`
	TotalCost     float64 `gorm:"type:decimal(10,2);comment:总成本" json:"total_cost"`

	// 供应商信息
	Supplier string `gorm:"type:varchar(200);comment:供应商" json:"supplier"`

	// 库存信息
	Stock int `gorm:"default:0;comment:库存数量" json:"stock"`

	// 状态
	IsActive bool `gorm:"default:true;comment:是否有效" json:"is_active"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (ProductCost) TableName() string {
	return "product_cost"
}

// BeforeCreate 创建前计算总成本
func (pc *ProductCost) BeforeCreate(tx *gorm.DB) error {
	// 计算总成本
	pc.TotalCost = pc.PurchasePrice + pc.ShippingCost + pc.HandlingCost + pc.OtherCost

	return nil
}

// BeforeUpdate 更新前重新计算总成本
func (pc *ProductCost) BeforeUpdate(tx *gorm.DB) error {
	pc.TotalCost = pc.PurchasePrice + pc.ShippingCost + pc.HandlingCost + pc.OtherCost
	return nil
}
