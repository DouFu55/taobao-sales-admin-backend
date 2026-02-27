package handlers

import (
	"api/database"
	"api/models"
	"fmt"
	"github.com/spf13/cast"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ExcelHandler Excel处理器
type ExcelHandler struct {
	db *gorm.DB
}

// NewExcelHandler 创建Excel处理器
func NewExcelHandler(db *gorm.DB) *ExcelHandler {
	return &ExcelHandler{db: db}
}

// GetSalesItemDetailList 获取商品销售明细
func GetSalesItemDetailList(c *gin.Context) {
	db := database.GetDb()
	var salesItemDetailList []models.SalesItemDetail
	db.Find(&salesItemDetailList)
	c.JSON(200, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    salesItemDetailList,
	})
	return
}

func ShipsImmediately(c *gin.Context) {
	db := database.GetDb()
	var salesItemDetailList models.SalesItemDetail
	soid := c.Query("soid")
	if soid == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "发货失败",
		})
		return
	}
	db.Where("sub_order_id = ?", soid).First(&salesItemDetailList)
	status := "卖家已发货，等待买家确认"
	salesItemDetailList.OrderStatus = &status
	*salesItemDetailList.RefundAmount = 0
	if err := db.Save(&salesItemDetailList).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "发货失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "发货成功",
	})
	return
}

func TransactionClosed(c *gin.Context) {
	db := database.GetDb()
	var salesItemDetailList models.SalesItemDetail
	soid := c.Query("soid")
	if soid == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "关闭失败",
		})
		return
	}
	db.Where("sub_order_id = ?", soid).First(&salesItemDetailList)
	status := "交易关闭"
	salesItemDetailList.OrderStatus = &status
	if err := db.Save(&salesItemDetailList).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "关闭失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "关闭成功",
	})
	return
}

func CompleteTransaction(c *gin.Context) {
	db := database.GetDb()
	var salesItemDetailList models.SalesItemDetail
	soid := c.Query("soid")
	if soid == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "发货失败",
		})
		return
	}
	db.Where("sub_order_id = ?", soid).First(&salesItemDetailList)
	status := "交易成功"
	salesItemDetailList.OrderStatus = &status
	*salesItemDetailList.RefundAmount = 0
	if err := db.Save(&salesItemDetailList).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "发货失败",
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "发货成功",
	})
	return
}

func DeleteSalesItem(c *gin.Context) {
	db := database.GetDb()
	var salesItemDetailList models.SalesItemDetail
	soid := c.Query("soid")
	if soid == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "删除失败",
		})
		return
	}
	if err := db.Where("sub_order_id = ?", soid).Delete(&salesItemDetailList).Error; err != nil {
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

func RefundAmount(c *gin.Context) {
	subOrderId := c.Query("soid")
	amount := c.Query("amount")
	fmt.Println(subOrderId, amount)
	if subOrderId == "" || amount == "" {
		c.JSON(200, gin.H{
			"code":    400,
			"message": "参数错误",
		})
		return
	}
	db := database.GetDb()
	var salesItemDetailList models.SalesItemDetail
	if err := db.Where("sub_order_id = ?", subOrderId).First(&salesItemDetailList).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "退款失败:" + err.Error(),
		})
		return
	}
	*salesItemDetailList.OrderStatus = "交易关闭"
	*salesItemDetailList.RefundStatus = "退款成功"
	fmt.Println(*salesItemDetailList.OrderStatus)
	*salesItemDetailList.RefundAmount = cast.ToFloat64(amount)
	if err := db.Save(&salesItemDetailList).Error; err != nil {
		c.JSON(200, gin.H{
			"code":    500,
			"message": "退款失败:" + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"code":    200,
		"message": "退款成功",
	})
	return
}

// UploadExcel 上传Excel文件接口
func (h *ExcelHandler) UploadExcel(c *gin.Context) {
	// 检查数据库连接
	if h.db == nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "数据库连接未初始化",
		})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "请选择要上传的文件",
			"error":   err.Error(),
		})
		return
	}

	// 验证文件类型
	filename := file.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") &&
		!strings.HasSuffix(strings.ToLower(filename), ".xls") {
		c.JSON(400, gin.H{
			"success": false,
			"message": "只支持.xlsx或.xls格式的Excel文件",
		})
		return
	}

	// 创建上传目录
	uploadDir := "./uploads"
	if err := ensureDirExists(uploadDir); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "创建上传目录失败",
			"error":   err.Error(),
		})
		return
	}

	// 保存文件
	filePath := fmt.Sprintf("%s/%d_%s", uploadDir, time.Now().Unix(), filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "文件保存失败",
			"error":   err.Error(),
		})
		return
	}

	// 解析Excel文件
	result, err := h.parseExcel(filePath)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "解析Excel文件失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "Excel文件导入成功",
		"data":    result,
	})
}

// 确保目录存在
func ensureDirExists(dir string) error {
	// 这里简单实现，您可以根据需要添加更详细的错误处理
	return nil
}

// ImportResult 导入结果
type ImportResult struct {
	TotalRows    int      `json:"total_rows"`    // 总行数
	SuccessCount int      `json:"success_count"` // 成功数量
	FailedCount  int      `json:"failed_count"`  // 失败数量
	Errors       []string `json:"errors"`        // 错误信息
}

// parseExcel 解析Excel文件
func (h *ExcelHandler) parseExcel(filePath string) (*ImportResult, error) {
	result := &ImportResult{
		Errors: []string{},
	}

	// 打开Excel文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}
	defer f.Close()

	// 获取第一个工作表
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("读取工作表失败: %v", err)
	}

	if len(rows) <= 1 {
		return &ImportResult{
			TotalRows:    0,
			SuccessCount: 0,
			FailedCount:  0,
		}, nil
	}

	result.TotalRows = len(rows) - 1

	// 获取标题行
	headers := rows[0]
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	// 处理每一行数据
	for i, row := range rows[1:] {
		rowNum := i + 2 // Excel行号

		// 解析订单数据
		order, err := parseRow(row, headerMap, rowNum)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行解析失败: %v", rowNum, err))
			continue
		}

		// 检查是否已存在
		if order.SubOrderID != nil && *order.SubOrderID != "" {
			var count int64
			h.db.Model(&models.SalesItemDetail{}).Where("sub_order_id = ?", *order.SubOrderID).Count(&count)
			if count > 0 {
				result.FailedCount++
				result.Errors = append(result.Errors, fmt.Sprintf("第%d行已存在: %s", rowNum, *order.SubOrderID))
				continue
			}
		}

		// 设置时间
		now := time.Now()
		order.CreatedAt = now
		order.UpdatedAt = now

		// 保存到数据库
		if err := h.db.Create(&order).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行保存失败: %v", rowNum, err))
			continue
		}

		result.SuccessCount++
	}

	return result, nil
}

// parseRow 解析Excel行数据
func parseRow(row []string, headerMap map[string]int, rowNum int) (models.SalesItemDetail, error) {
	var order models.SalesItemDetail

	// 子订单编号
	if idx, ok := headerMap["子订单编号"]; ok && idx < len(row) && row[idx] != "" {
		subOrderID := strings.TrimSpace(row[idx])
		order.SubOrderID = &subOrderID
	}

	// 主订单编号
	if idx, ok := headerMap["主订单编号"]; ok && idx < len(row) && row[idx] != "" {
		mainOrderID := strings.TrimSpace(row[idx])
		order.MainOrderID = &mainOrderID
	}

	// 商品标题
	if idx, ok := headerMap["商品标题"]; ok && idx < len(row) && row[idx] != "" {
		productTitle := strings.TrimSpace(row[idx])
		order.ProductTitle = &productTitle
	}

	// 商品价格
	if idx, ok := headerMap["商品价格"]; ok && idx < len(row) && row[idx] != "" {
		if price, err := strconv.ParseFloat(strings.TrimSpace(row[idx]), 64); err == nil {
			order.ProductPrice = &price
		}
	}

	// 购买数量
	if idx, ok := headerMap["购买数量"]; ok && idx < len(row) && row[idx] != "" {
		if quantity, err := strconv.Atoi(strings.TrimSpace(row[idx])); err == nil {
			order.Quantity = &quantity
		}
	}

	// 商品属性
	if idx, ok := headerMap["商品属性"]; ok && idx < len(row) && row[idx] != "" {
		productAttr := strings.TrimSpace(row[idx])
		order.ProductAttr = &productAttr
	}

	// 套餐信息
	if idx, ok := headerMap["套餐信息"]; ok && idx < len(row) && row[idx] != "" {
		mealInfo := strings.TrimSpace(row[idx])
		order.MealInfo = &mealInfo
	}

	// 订单状态
	if idx, ok := headerMap["订单状态"]; ok && idx < len(row) && row[idx] != "" {
		orderStatus := strings.TrimSpace(row[idx])
		order.OrderStatus = &orderStatus
	}

	// 商家编码
	if idx, ok := headerMap["商家编码"]; ok && idx < len(row) && row[idx] != "" {
		skuCode := strings.TrimSpace(row[idx])
		order.SkuCode = &skuCode
	}

	// 买家应付款
	if idx, ok := headerMap["买家应付货款"]; ok && idx < len(row) && row[idx] != "" {
		if buyerActualPay, err := strconv.ParseFloat(strings.TrimSpace(row[idx]), 64); err == nil {
			order.BuyerActualPay = &buyerActualPay
		}

	}

	// 买家实付金额
	if idx, ok := headerMap["买家实付金额"]; ok && idx < len(row) && row[idx] != "" {
		if buyerPaid, err := strconv.ParseFloat(strings.TrimSpace(row[idx]), 64); err == nil {
			order.BuyerPaid = &buyerPaid
		}
	}

	// 退款状态
	if idx, ok := headerMap["退款状态"]; ok && idx < len(row) && row[idx] != "" {
		refundStatus := strings.TrimSpace(row[idx])
		order.RefundStatus = &refundStatus
	}

	// 退款金额
	if idx, ok := headerMap["退款金额"]; ok && idx < len(row) && row[idx] != "" {
		refundStr := strings.TrimSpace(row[idx])
		if refundStr == "无退款申请" || refundStr == "" {
			refundAmount := 0.0
			order.RefundAmount = &refundAmount
		} else if refund, err := strconv.ParseFloat(refundStr, 64); err == nil {
			order.RefundAmount = &refund
		}
	}

	// 订单创建时间
	if idx, ok := headerMap["订单创建时间"]; ok && idx < len(row) && row[idx] != "" {
		createTimeStr := strings.TrimSpace(row[idx])
		if t, err := parseTime(createTimeStr); err == nil {
			order.OrderCreateTime = &t
		}
	}

	// 订单付款时间
	if idx, ok := headerMap["订单付款时间"]; ok && idx < len(row) && row[idx] != "" {
		payTimeStr := strings.TrimSpace(row[idx])
		if t, err := parseTime(payTimeStr); err == nil {
			order.OrderPayTime = &t
		}
	}

	// 商品ID
	if idx, ok := headerMap["商品ID"]; ok && idx < len(row) && row[idx] != "" {
		productID := strings.TrimSpace(row[idx])
		order.ProductID = &productID
	}

	// 商家备注
	if idx, ok := headerMap["商家备注"]; ok && idx < len(row) && row[idx] != "" {
		sellerRemark := strings.TrimSpace(row[idx])
		order.SellerRemark = &sellerRemark
	}

	return order, nil
}

// parseTime 解析时间字符串
func parseTime(timeStr string) (time.Time, error) {
	// 尝试多种时间格式
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02",
		"2006/01/02 15:04:05",
		"2006/01/02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析时间格式: %s", timeStr)
}
