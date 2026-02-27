package handlers

import (
	"api/database"
	"api/models"
	"api/validators"
	"github.com/gin-gonic/gin"
	"time"
)

func GetPromotionList(c *gin.Context) {
	db := database.GetDb()
	var productCost []models.PromotionExpense
	db.Find(&productCost)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    productCost,
	})
	return
}

func CreatePromotionItem(c *gin.Context) {
	var r validators.PromotionExpenseValidator
	err := c.ShouldBind(&r)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 将字符串日期转换为 time.Time
	promotionDate, err := time.Parse("2006-01-02", r.PromotionDate)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "日期格式错误，请使用 YYYY-MM-DD 格式",
		})
		return
	}

	var promotionItem models.PromotionExpense
	promotionItem.PromotionDate = promotionDate
	promotionItem.PromotionType = r.PromotionType
	promotionItem.PlanID = r.PlanID
	promotionItem.ExpenseAmount = r.ExpenseAmount
	promotionItem.OrdersCount = r.OrdersCount
	promotionItem.Remark = r.Remark

	var db = database.GetDb()
	if db.Create(&promotionItem).Error != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "创建失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

func EditPromotionItem(c *gin.Context) {
	var r validators.PromotionExpenseValidator
	err := c.ShouldBind(&r)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}
	if r.Id == 0 {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}

	// 查询数据
	var db = database.GetDb()
	var promotionItem models.PromotionExpense
	if err = db.Where("id = ?", r.Id).First(&promotionItem).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "编辑失败",
		})
		return
	}

	// 将字符串日期转换为 time.Time
	promotionDate, err := time.Parse("2006-01-02", r.PromotionDate)
	if err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "日期格式错误，请使用 YYYY-MM-DD 格式",
		})
		return
	}

	promotionItem.PromotionDate = promotionDate
	promotionItem.PromotionType = r.PromotionType
	promotionItem.PlanID = r.PlanID
	promotionItem.ExpenseAmount = r.ExpenseAmount
	promotionItem.OrdersCount = r.OrdersCount
	promotionItem.Remark = r.Remark

	if db.Save(&promotionItem).Error != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "编辑失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "编辑成功",
	})
}

func DeletePromotionItem(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	var db = database.GetDb()
	var promotionItem models.PromotionExpense
	if err := db.Where("id = ?", id).Delete(&promotionItem).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "删除失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "删除成功",
	})
	return
}
