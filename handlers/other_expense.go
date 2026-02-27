package handlers

import (
	"api/database"
	"api/models"
	"github.com/gin-gonic/gin"
)

func GetOtherExpenseList(c *gin.Context) {
	db := database.GetDb()
	var otherExpenseList []models.OtherExpense
	if err := db.Find(&otherExpenseList).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "获取失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    otherExpenseList,
	})
}

func CreateOtherExpenseItem(c *gin.Context) {
	var r models.OtherExpense
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}
	var otherExpenseItem models.OtherExpense
	otherExpenseItem.ExpenseDate = r.ExpenseDate
	otherExpenseItem.ExpenseType = r.ExpenseType
	otherExpenseItem.ItemName = r.ItemName
	otherExpenseItem.Amount = r.Amount
	otherExpenseItem.SubOrderNo = r.SubOrderNo
	otherExpenseItem.Remark = r.Remark
	db := database.GetDb()
	if err := db.Create(&otherExpenseItem).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "创建失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

func EditOtherExpenseItem(c *gin.Context) {
	var r models.OtherExpense
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误：" + err.Error(),
		})
		return
	}

	db := database.GetDb()
	var otherExpenseItem models.OtherExpense
	// 查询数据
	if err := db.Where("id = ?", r.ID).First(&otherExpenseItem).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "编辑失败",
		})
		return
	}
	otherExpenseItem.ExpenseDate = r.ExpenseDate
	otherExpenseItem.ExpenseType = r.ExpenseType
	otherExpenseItem.ItemName = r.ItemName
	otherExpenseItem.Amount = r.Amount
	otherExpenseItem.SubOrderNo = r.SubOrderNo
	otherExpenseItem.Remark = r.Remark
	if err := db.Save(&otherExpenseItem).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "创建失败：" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "创建成功",
	})
}

func DeleteOtherExpenseItem(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	var db = database.GetDb()
	var otherExpense models.OtherExpense
	if err := db.Where("id = ?", id).Delete(&otherExpense).Error; err != nil {
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
