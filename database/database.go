package database

import (
	"api/models"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() error {
	// 使用纯 Go 的 SQLite 驱动
	db, err := gorm.Open(sqlite.Open("database/tb.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("创建数据库连接失败: %v", err)
	}

	// 自动迁移
	err = db.AutoMigrate(&models.SalesItemDetail{})
	err = db.AutoMigrate(&models.ProductCost{})
	err = db.AutoMigrate(&models.PromotionExpense{})
	err = db.AutoMigrate(&models.OtherExpense{})
	if err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	DB = db

	return nil
}

func GetDb() *gorm.DB { return DB }
