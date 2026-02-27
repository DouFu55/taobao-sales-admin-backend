package handlers

import (
	"api/database"
	"api/models"
	"api/validators"
	"fmt"
	"github.com/gin-gonic/gin"
)

func ProductCostList(c *gin.Context) {
	db := database.GetDb()
	var productCost []models.ProductCost
	db.Find(&productCost)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    productCost,
	})
	return
}

func CreateProductCost(c *gin.Context) {
	var r validators.ProductCostValidator
	err := c.ShouldBind(&r) // 这里需要传递指针 &r
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	var productCost models.ProductCost
	productCost.SKU = r.SKU
	productCost.ProductID = uint64(r.ProductID)

	// 处理指针字段，避免空指针解引用
	if r.ProductName != "" {
		productCost.ProductName = r.ProductName
	} else {
		productCost.ProductName = "" // 默认值
	}

	if r.SpecInfo != "" {
		productCost.SpecInfo = r.SpecInfo
	} else {
		productCost.SpecInfo = "" // 默认值
	}

	if r.PurchasePrice != nil {
		productCost.PurchasePrice = *r.PurchasePrice
	} else {
		productCost.PurchasePrice = 0 // 默认值
	}

	if r.ShippingCost != nil {
		productCost.ShippingCost = *r.ShippingCost
	} else {
		productCost.ShippingCost = 0 // 默认值
	}

	if r.HandlingCost != nil {
		productCost.HandlingCost = *r.HandlingCost
	} else {
		productCost.HandlingCost = 0 // 默认值
	}

	if r.OtherCost != nil {
		productCost.OtherCost = *r.OtherCost
	} else {
		productCost.OtherCost = 0 // 默认值
	}

	// 计算总成本（如果前端没传）
	if r.TotalCost != nil {
		productCost.TotalCost = *r.TotalCost
	} else {
		productCost.TotalCost = productCost.PurchasePrice +
			productCost.ShippingCost +
			productCost.HandlingCost +
			productCost.OtherCost
	}

	if r.Supplier != "" {
		productCost.Supplier = r.Supplier
	} else {
		productCost.Supplier = "" // 默认值
	}

	if r.Stock != nil {
		productCost.Stock = *r.Stock
	} else {
		productCost.Stock = 0 // 默认值
	}

	if r.IsActive != nil {
		productCost.IsActive = *r.IsActive
	} else {
		productCost.IsActive = true // 默认值
	}

	db := database.GetDb()
	if err := db.Create(&productCost).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
	})
	return
}

func DeleteCostItem(c *gin.Context) {
	var id = c.Query("id")
	if id == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	db := database.GetDb()
	var productCost models.ProductCost
	if err := db.Where("id = ?", id).Delete(&productCost).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "删除失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
	return
}

func EditCostItem(c *gin.Context) {
	var r validators.ProductCostValidator
	err := c.ShouldBind(&r) // 这里需要传递指针 &r
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	var productCost models.ProductCost
	// 查询数据
	db := database.GetDb()
	fmt.Println(r)
	if err = db.Where("id = ?", r.ID).First(&productCost).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "编辑失败",
		})
		return
	}

	productCost.SKU = r.SKU
	productCost.ProductID = uint64(r.ProductID)

	// 处理指针字段，避免空指针解引用
	if r.ProductName != "" {
		productCost.ProductName = r.ProductName
	} else {
		productCost.ProductName = "" // 默认值
	}

	if r.SpecInfo != "" {
		productCost.SpecInfo = r.SpecInfo
	} else {
		productCost.SpecInfo = "" // 默认值
	}

	if r.PurchasePrice != nil {
		productCost.PurchasePrice = *r.PurchasePrice
	} else {
		productCost.PurchasePrice = 0 // 默认值
	}

	if r.ShippingCost != nil {
		productCost.ShippingCost = *r.ShippingCost
	} else {
		productCost.ShippingCost = 0 // 默认值
	}

	if r.HandlingCost != nil {
		productCost.HandlingCost = *r.HandlingCost
	} else {
		productCost.HandlingCost = 0 // 默认值
	}

	if r.OtherCost != nil {
		productCost.OtherCost = *r.OtherCost
	} else {
		productCost.OtherCost = 0 // 默认值
	}

	// 计算总成本（如果前端没传）
	if r.TotalCost != nil {
		productCost.TotalCost = *r.TotalCost
	} else {
		productCost.TotalCost = productCost.PurchasePrice +
			productCost.ShippingCost +
			productCost.HandlingCost +
			productCost.OtherCost
	}

	if r.Supplier != "" {
		productCost.Supplier = r.Supplier
	} else {
		productCost.Supplier = "" // 默认值
	}

	if r.Stock != nil {
		productCost.Stock = *r.Stock
	} else {
		productCost.Stock = 0 // 默认值
	}

	if r.IsActive != nil {
		productCost.IsActive = *r.IsActive
	} else {
		productCost.IsActive = true // 默认值
	}
	if err = db.Save(&productCost).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
	})
	return
}
