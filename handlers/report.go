package handlers

import (
	"api/database"
	"api/models"
	"api/types"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

func GetDailyReport(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
	}

	startOfDay, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(200, gin.H{"code": 400, "message": "日期格式错误"})
		return
	}
	endOfDay := startOfDay.Add(24 * time.Hour)

	var dailyReport types.DailyReport
	dailyReport.ReportDate = dateStr
	db := database.GetDb()

	var salesItems []models.SalesItemDetail
	var promotionExpenses []models.PromotionExpense
	var otherExpenses []models.OtherExpense
	var productCosts []models.ProductCost

	db.Where("order_pay_time >= ? AND order_pay_time < ? AND buyer_actual_pay > 0", startOfDay, endOfDay).Find(&salesItems)
	db.Where("promotion_date >= ? AND promotion_date < ?", startOfDay, endOfDay).Find(&promotionExpenses)
	db.Where("expense_date >= ? AND expense_date < ?", startOfDay, endOfDay).Find(&otherExpenses)
	db.Find(&productCosts)

	skuIndex := make(map[string]models.ProductCost)
	prodSpecIndex := make(map[string]models.ProductCost)
	for _, pc := range productCosts {
		if !pc.IsActive {
			continue
		}
		// 统一去除空格，确保匹配精准
		sKey := strings.TrimSpace(pc.SKU)
		if sKey != "" {
			skuIndex[sKey] = pc
		}
		if pc.ProductID > 0 && pc.SpecInfo != "" {
			// 构建索引：使用数字ID + 清洗后的规格字符串
			key := fmt.Sprintf("%d|%s", pc.ProductID, strings.TrimSpace(pc.SpecInfo))
			prodSpecIndex[key] = pc
		}
	}

	orderMap := make(map[string]bool)
	for _, item := range salesItems {
		paid := item.GetBuyerActualPay()
		qty := 0
		if item.Quantity != nil {
			qty = *item.Quantity
		}

		dailyReport.SalesSummary.TotalAmount += paid
		dailyReport.SalesSummary.NetAmount += paid
		dailyReport.SalesSummary.ItemCount += qty

		if item.MainOrderID != nil {
			orderMap[*item.MainOrderID] = true
		}

		if item.RefundAmount != nil && *item.RefundAmount > 0 {
			dailyReport.SalesSummary.RefundAmount += *item.RefundAmount
			dailyReport.SalesSummary.RefundCount++
		}

		status := ""
		if item.OrderStatus != nil {
			status = *item.OrderStatus
		}
		if status == "卖家已发货，等待买家确认" || status == "交易成功" {
			// 调用修复后的 findCost
			costDetail := findCost(item, skuIndex, prodSpecIndex)
			if costDetail.Found {
				qf := float64(qty)
				dailyReport.CostSummary.ProductCost += costDetail.PurchasePrice * qf
				dailyReport.CostSummary.ShippingCost += costDetail.ShippingCost * qf
				dailyReport.CostSummary.HandlingCost += costDetail.HandlingCost * qf
				dailyReport.CostSummary.OtherCost += costDetail.OtherCost * qf
				dailyReport.CostSummary.TotalCost += costDetail.TotalCost * qf
			}
		}
	}

	dailyReport.SalesSummary.OrderCount = len(orderMap)
	dailyReport.SalesSummary.NetAmount -= dailyReport.SalesSummary.RefundAmount

	if dailyReport.SalesSummary.OrderCount > 0 {
		dailyReport.SalesSummary.AvgOrderValue = dailyReport.SalesSummary.TotalAmount / float64(dailyReport.SalesSummary.OrderCount)
	}

	for _, pe := range promotionExpenses {
		dailyReport.ExpenseSummary.PromotionExpenses += pe.ExpenseAmount
	}
	for _, oe := range otherExpenses {
		dailyReport.ExpenseSummary.OtherExpenses += oe.Amount
	}
	dailyReport.ExpenseSummary.TotalExpenses = dailyReport.ExpenseSummary.PromotionExpenses + dailyReport.ExpenseSummary.OtherExpenses

	dailyReport.ProfitAnalysis.GrossProfit = dailyReport.SalesSummary.NetAmount - dailyReport.CostSummary.TotalCost
	dailyReport.ProfitAnalysis.OperatingProfit = dailyReport.ProfitAnalysis.GrossProfit - dailyReport.ExpenseSummary.TotalExpenses
	dailyReport.ProfitAnalysis.NetProfit = dailyReport.ProfitAnalysis.OperatingProfit

	if dailyReport.SalesSummary.NetAmount > 0 {
		gMargin := (dailyReport.ProfitAnalysis.GrossProfit / dailyReport.SalesSummary.NetAmount) * 100
		nMargin := (dailyReport.ProfitAnalysis.NetProfit / dailyReport.SalesSummary.NetAmount) * 100
		dailyReport.ProfitAnalysis.GrossMargin = fmt.Sprintf("%.2f%%", gMargin)
		dailyReport.ProfitAnalysis.NetMargin = fmt.Sprintf("%.2f%%", nMargin)
	} else {
		dailyReport.ProfitAnalysis.GrossMargin = "0.00%"
		dailyReport.ProfitAnalysis.NetMargin = "0.00%"
	}

	c.JSON(200, gin.H{"code": 200, "message": "获取成功", "data": dailyReport})
}

func GetRangeReport(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")
	if startDateStr == "" || endDateStr == "" {
		c.JSON(200, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	startTime, err1 := parseFlexibleDate(startDateStr)
	endTime, err2 := parseFlexibleDate(endDateStr)
	if err1 != nil || err2 != nil {
		c.JSON(200, gin.H{"code": 400, "message": "日期格式错误"})
		return
	}

	realQueryEndTime := endTime
	if len(endDateStr) <= 10 {
		realQueryEndTime = endTime.Add(24 * time.Hour).Add(-time.Second)
	}

	db := database.GetDb()
	var salesItems []models.SalesItemDetail
	var promoExpenses []models.PromotionExpense
	var otherExpenses []models.OtherExpense
	var productCosts []models.ProductCost
	var report types.DailyReport
	report.ReportDate = fmt.Sprintf("%s 至 %s", startDateStr, endDateStr)

	db.Where("order_pay_time >= ? AND order_pay_time <= ? AND buyer_actual_pay > 0", startTime, realQueryEndTime).Find(&salesItems)
	db.Where("promotion_date >= ? AND promotion_date <= ?", startTime, realQueryEndTime).Find(&promoExpenses)
	db.Where("expense_date >= ? AND expense_date <= ?", startTime, realQueryEndTime).Find(&otherExpenses)
	db.Find(&productCosts)

	skuMap := make(map[string]models.ProductCost)
	specMap := make(map[string]models.ProductCost)
	for _, pc := range productCosts {
		if !pc.IsActive {
			continue
		}
		sKey := strings.TrimSpace(pc.SKU)
		if sKey != "" {
			skuMap[sKey] = pc
		}
		if pc.ProductID > 0 && pc.SpecInfo != "" {
			key := fmt.Sprintf("%d|%s", pc.ProductID, strings.TrimSpace(pc.SpecInfo))
			specMap[key] = pc
		}
	}

	orderMap := make(map[string]bool)
	for _, item := range salesItems {
		paid := item.GetBuyerActualPay()
		qty := 0
		if item.Quantity != nil {
			qty = *item.Quantity
		}
		if item.MainOrderID != nil {
			orderMap[*item.MainOrderID] = true
		}

		report.SalesSummary.TotalAmount += paid
		report.SalesSummary.NetAmount += paid
		report.SalesSummary.ItemCount += qty

		if item.RefundAmount != nil && *item.RefundAmount > 0 {
			report.SalesSummary.RefundAmount += *item.RefundAmount
			report.SalesSummary.RefundCount++
		}

		status := ""
		if item.OrderStatus != nil {
			status = *item.OrderStatus
		}
		if status == "卖家已发货，等待买家确认" || status == "交易成功" {
			costDetail := findCost(item, skuMap, specMap)
			if costDetail.Found {
				qf := float64(qty)
				report.CostSummary.ProductCost += costDetail.PurchasePrice * qf
				report.CostSummary.ShippingCost += costDetail.ShippingCost * qf
				report.CostSummary.HandlingCost += costDetail.HandlingCost * qf
				report.CostSummary.OtherCost += costDetail.OtherCost * qf
				report.CostSummary.TotalCost += costDetail.TotalCost * qf
			}
		}
	}

	report.SalesSummary.OrderCount = len(orderMap)
	report.SalesSummary.NetAmount -= report.SalesSummary.RefundAmount

	if report.SalesSummary.OrderCount > 0 {
		report.SalesSummary.AvgOrderValue = report.SalesSummary.TotalAmount / float64(report.SalesSummary.OrderCount)
	}

	for _, pe := range promoExpenses {
		report.ExpenseSummary.PromotionExpenses += pe.ExpenseAmount
	}
	for _, oe := range otherExpenses {
		report.ExpenseSummary.OtherExpenses += oe.Amount
	}
	report.ExpenseSummary.TotalExpenses = report.ExpenseSummary.PromotionExpenses + report.ExpenseSummary.OtherExpenses

	report.ProfitAnalysis.GrossProfit = report.SalesSummary.NetAmount - report.CostSummary.TotalCost
	report.ProfitAnalysis.OperatingProfit = report.ProfitAnalysis.GrossProfit - report.ExpenseSummary.TotalExpenses
	report.ProfitAnalysis.NetProfit = report.ProfitAnalysis.OperatingProfit

	if report.SalesSummary.NetAmount > 0 {
		gm := (report.ProfitAnalysis.GrossProfit / report.SalesSummary.NetAmount) * 100
		nm := (report.ProfitAnalysis.NetProfit / report.SalesSummary.NetAmount) * 100
		report.ProfitAnalysis.GrossMargin = fmt.Sprintf("%.2f%%", gm)
		report.ProfitAnalysis.NetMargin = fmt.Sprintf("%.2f%%", nm)
	} else {
		report.ProfitAnalysis.GrossMargin = "0.00%"
		report.ProfitAnalysis.NetMargin = "0.00%"
	}

	c.JSON(200, gin.H{"code": 200, "message": "获取成功", "data": report})
}

func GetMonthlyReport(c *gin.Context) {
	monthStr := c.Query("month")
	if monthStr == "" {
		monthStr = time.Now().Format("2006-01")
	}
	t, err := time.Parse("2006-01", monthStr)
	if err != nil {
		t, _ = time.Parse("2006-1", monthStr)
	}
	startOfMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	var report types.DailyReport
	report.ReportDate = startOfMonth.Format("2006年01月")
	db := database.GetDb()

	var salesItems []models.SalesItemDetail
	var promotionExpenses []models.PromotionExpense
	var otherExpenses []models.OtherExpense
	var productCosts []models.ProductCost

	db.Where("order_pay_time >= ? AND order_pay_time < ? AND buyer_actual_pay > 0", startOfMonth, endOfMonth).Find(&salesItems)
	db.Where("promotion_date >= ? AND promotion_date < ?", startOfMonth, endOfMonth).Find(&promotionExpenses)
	db.Where("expense_date >= ? AND expense_date < ?", startOfMonth, endOfMonth).Find(&otherExpenses)
	db.Find(&productCosts)

	skuIndex := make(map[string]models.ProductCost)
	prodSpecIndex := make(map[string]models.ProductCost)
	for _, pc := range productCosts {
		if !pc.IsActive {
			continue
		}
		sKey := strings.TrimSpace(pc.SKU)
		if sKey != "" {
			skuIndex[sKey] = pc
		}
		if pc.ProductID > 0 && pc.SpecInfo != "" {
			key := fmt.Sprintf("%d|%s", pc.ProductID, strings.TrimSpace(pc.SpecInfo))
			prodSpecIndex[key] = pc
		}
	}

	orderMap := make(map[string]bool)
	for _, item := range salesItems {
		paid := item.GetBuyerActualPay()
		qty := 0
		if item.Quantity != nil {
			qty = *item.Quantity
		}
		report.SalesSummary.TotalAmount += paid
		report.SalesSummary.NetAmount += paid
		report.SalesSummary.ItemCount += qty
		if item.MainOrderID != nil {
			orderMap[*item.MainOrderID] = true
		}

		if item.RefundAmount != nil && *item.RefundAmount > 0 {
			report.SalesSummary.RefundAmount += *item.RefundAmount
			report.SalesSummary.RefundCount++
		}

		status := ""
		if item.OrderStatus != nil {
			status = *item.OrderStatus
		}
		if status == "卖家已发货，等待买家确认" || status == "交易成功" {
			costDetail := findCost(item, skuIndex, prodSpecIndex)
			if costDetail.Found {
				qf := float64(qty)
				report.CostSummary.ProductCost += costDetail.PurchasePrice * qf
				report.CostSummary.ShippingCost += costDetail.ShippingCost * qf
				report.CostSummary.HandlingCost += costDetail.HandlingCost * qf
				report.CostSummary.OtherCost += costDetail.OtherCost * qf
				report.CostSummary.TotalCost += costDetail.TotalCost * qf
			}
		}
	}
	report.SalesSummary.OrderCount = len(orderMap)
	report.SalesSummary.NetAmount -= report.SalesSummary.RefundAmount
	if report.SalesSummary.OrderCount > 0 {
		report.SalesSummary.AvgOrderValue = report.SalesSummary.TotalAmount / float64(report.SalesSummary.OrderCount)
	}
	for _, pe := range promotionExpenses {
		report.ExpenseSummary.PromotionExpenses += pe.ExpenseAmount
	}
	for _, oe := range otherExpenses {
		report.ExpenseSummary.OtherExpenses += oe.Amount
	}
	report.ExpenseSummary.TotalExpenses = report.ExpenseSummary.PromotionExpenses + report.ExpenseSummary.OtherExpenses
	report.ProfitAnalysis.GrossProfit = report.SalesSummary.NetAmount - report.CostSummary.TotalCost
	report.ProfitAnalysis.OperatingProfit = report.ProfitAnalysis.GrossProfit - report.ExpenseSummary.TotalExpenses
	report.ProfitAnalysis.NetProfit = report.ProfitAnalysis.OperatingProfit
	if report.SalesSummary.NetAmount > 0 {
		gm := (report.ProfitAnalysis.GrossProfit / report.SalesSummary.NetAmount) * 100
		nm := (report.ProfitAnalysis.NetProfit / report.SalesSummary.NetAmount) * 100
		report.ProfitAnalysis.GrossMargin = fmt.Sprintf("%.2f%%", gm)
		report.ProfitAnalysis.NetMargin = fmt.Sprintf("%.2f%%", nm)
	}
	c.JSON(200, gin.H{"code": 200, "message": "获取成功", "data": report})
}

func GetYearlyReport(c *gin.Context) {
	yearStr := c.Query("year")
	if yearStr == "" {
		yearStr = time.Now().Format("2006")
	}
	t, _ := time.Parse("2006", yearStr)
	startOfYear := time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
	endOfYear := startOfYear.AddDate(1, 0, 0)

	var report types.DailyReport
	report.ReportDate = fmt.Sprintf("%d年度", t.Year())
	db := database.GetDb()
	var salesItems []models.SalesItemDetail
	var promotionExpenses []models.PromotionExpense
	var otherExpenses []models.OtherExpense
	var productCosts []models.ProductCost

	db.Where("order_pay_time >= ? AND order_pay_time < ? AND buyer_actual_pay > 0", startOfYear, endOfYear).Find(&salesItems)
	db.Where("promotion_date >= ? AND promotion_date < ?", startOfYear, endOfYear).Find(&promotionExpenses)
	db.Where("expense_date >= ? AND expense_date < ?", startOfYear, endOfYear).Find(&otherExpenses)
	db.Find(&productCosts)

	skuIndex := make(map[string]models.ProductCost)
	prodSpecIndex := make(map[string]models.ProductCost)
	for _, pc := range productCosts {
		if !pc.IsActive {
			continue
		}
		sKey := strings.TrimSpace(pc.SKU)
		if sKey != "" {
			skuIndex[sKey] = pc
		}
		if pc.ProductID > 0 && pc.SpecInfo != "" {
			key := fmt.Sprintf("%d|%s", pc.ProductID, strings.TrimSpace(pc.SpecInfo))
			prodSpecIndex[key] = pc
		}
	}

	orderMap := make(map[string]bool)
	for _, item := range salesItems {
		paid := item.GetBuyerActualPay()
		qty := 0
		if item.Quantity != nil {
			qty = *item.Quantity
		}
		report.SalesSummary.TotalAmount += paid
		report.SalesSummary.NetAmount += paid
		report.SalesSummary.ItemCount += qty
		if item.MainOrderID != nil {
			orderMap[*item.MainOrderID] = true
		}

		if item.RefundAmount != nil && *item.RefundAmount > 0 {
			report.SalesSummary.RefundAmount += *item.RefundAmount
			report.SalesSummary.RefundCount++
		}

		status := ""
		if item.OrderStatus != nil {
			status = *item.OrderStatus
		}
		if status == "卖家已发货，等待买家确认" || status == "交易成功" {
			costDetail := findCost(item, skuIndex, prodSpecIndex)
			if costDetail.Found {
				qf := float64(qty)
				report.CostSummary.ProductCost += costDetail.PurchasePrice * qf
				report.CostSummary.ShippingCost += costDetail.ShippingCost * qf
				report.CostSummary.HandlingCost += costDetail.HandlingCost * qf
				report.CostSummary.OtherCost += costDetail.OtherCost * qf
				report.CostSummary.TotalCost += costDetail.TotalCost * qf
			}
		}
	}
	report.SalesSummary.OrderCount = len(orderMap)
	report.SalesSummary.NetAmount -= report.SalesSummary.RefundAmount
	if report.SalesSummary.OrderCount > 0 {
		report.SalesSummary.AvgOrderValue = report.SalesSummary.TotalAmount / float64(report.SalesSummary.OrderCount)
	}
	for _, pe := range promotionExpenses {
		report.ExpenseSummary.PromotionExpenses += pe.ExpenseAmount
	}
	for _, oe := range otherExpenses {
		report.ExpenseSummary.OtherExpenses += oe.Amount
	}
	report.ExpenseSummary.TotalExpenses = report.ExpenseSummary.PromotionExpenses + report.ExpenseSummary.OtherExpenses
	report.ProfitAnalysis.GrossProfit = report.SalesSummary.NetAmount - report.CostSummary.TotalCost
	report.ProfitAnalysis.OperatingProfit = report.ProfitAnalysis.GrossProfit - report.ExpenseSummary.TotalExpenses
	report.ProfitAnalysis.NetProfit = report.ProfitAnalysis.OperatingProfit
	if report.SalesSummary.NetAmount > 0 {
		gm := (report.ProfitAnalysis.GrossProfit / report.SalesSummary.NetAmount) * 100
		nm := (report.ProfitAnalysis.NetProfit / report.SalesSummary.NetAmount) * 100
		report.ProfitAnalysis.GrossMargin = fmt.Sprintf("%.2f%%", gm)
		report.ProfitAnalysis.NetMargin = fmt.Sprintf("%.2f%%", nm)
	}
	c.JSON(200, gin.H{"code": 200, "message": "获取成功", "data": report})
}

// 修复后的匹配逻辑
func findCost(sale models.SalesItemDetail, skuIndex map[string]models.ProductCost, prodSpecIndex map[string]models.ProductCost) types.ProductCostDetail {
	var matched *models.ProductCost

	// 1. 尝试使用 SkuCode 匹配 (商家编码)
	sku := ""
	if sale.SkuCode != nil {
		sku = strings.TrimSpace(*sale.SkuCode)
	}
	if sku != "" {
		if c, ok := skuIndex[sku]; ok {
			matched = &c
		}
	}

	// 2. 如果 SkuCode 为空或匹配失败，使用 ProductID + ProductAttr 匹配 (颜色分类)
	if matched == nil && sale.ProductID != nil && sale.ProductAttr != nil {
		// 将 string 类型的 ProductID 转为 uint64，以便与 ProductCost.ProductID 对应的 Key 格式对齐
		id, err := strconv.ParseUint(strings.TrimSpace(*sale.ProductID), 10, 64)
		if err == nil {
			// 生成 Key 逻辑必须与索引处 fmt.Sprintf("%d|%s", pc.ProductID, pc.SpecInfo) 完全一致
			key := fmt.Sprintf("%d|%s", id, strings.TrimSpace(*sale.ProductAttr))
			if c, ok := prodSpecIndex[key]; ok {
				matched = &c
			}
		}
	}

	if matched != nil {
		return types.ProductCostDetail{
			PurchasePrice: matched.PurchasePrice,
			ShippingCost:  matched.ShippingCost,
			HandlingCost:  matched.HandlingCost,
			OtherCost:     matched.OtherCost,
			TotalCost:     matched.TotalCost,
			Found:         true,
		}
	}
	return types.ProductCostDetail{Found: false}
}

func parseFlexibleDate(dateStr string) (time.Time, error) {
	if len(dateStr) <= 10 {
		return time.Parse("2006-01-02", dateStr)
	}
	return time.Parse("2006-01-02 15:04:05", dateStr)
}
