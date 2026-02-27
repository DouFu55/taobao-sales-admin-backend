package main

import (
	"api/database"
	"api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	db := database.GetDb()
	excelHandler := handlers.NewExcelHandler(db)

	// 只保留上传接口
	r.POST("/api/sales/execl/upload", excelHandler.UploadExcel)

	// 获取数据列表
	r.GET("/api/sales", handlers.GetSalesItemDetailList)
	// 订单发货
	r.GET("/api/sales/ship", handlers.ShipsImmediately)
	// 交易关闭
	r.GET("/api/sales/close", handlers.TransactionClosed)
	// 交易完成
	r.GET("/api/sales/complete", handlers.CompleteTransaction)
	r.DELETE("/api/sales/delete", handlers.DeleteSalesItem)
	r.GET("/api/sales/refund", handlers.RefundAmount)

	r.GET("/api/product/cost", handlers.ProductCostList)
	r.POST("/api/product/cost", handlers.CreateProductCost)
	r.DELETE("/api/product/cost", handlers.DeleteCostItem)
	r.PUT("/api/product/cost", handlers.EditCostItem)

	r.GET("/api/product/promotion", handlers.GetPromotionList)
	r.POST("/api/product/promotion", handlers.CreatePromotionItem)
	r.PUT("/api/product/promotion", handlers.EditPromotionItem)
	r.DELETE("/api/product/promotion", handlers.DeletePromotionItem)

	r.GET("/api/product/otherexpense", handlers.GetOtherExpenseList)
	r.POST("/api/product/otherexpense", handlers.CreateOtherExpenseItem)
	r.PUT("/api/product/otherexpense", handlers.EditOtherExpenseItem)
	r.DELETE("/api/product/otherexpense", handlers.DeleteOtherExpenseItem)

	//r.GET("/api/report/daily", handlers.GetDailyReport)

	// 基础日报接口 (传入参数 ?date=2026-01-20)
	r.GET("/api/report/daily", handlers.GetDailyReport)
	// 自定义日期范围接口 (传入参数 ?start_date=2026-01-25&end_date=2026-01-29)
	r.GET("/api/report/range", handlers.GetRangeReport)
	// 月度报表接口 (传入参数 ?month=2026-01)
	r.GET("/api/report/monthly", handlers.GetMonthlyReport)
	// 年度报表接口 (传入参数 ?year=2026)
	r.GET("/api/report/yearly", handlers.GetYearlyReport)
}
