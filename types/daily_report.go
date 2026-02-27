package types

// DailyReport 日报数据结构
type DailyReport struct {
	ReportDate     string         `json:"report_date"`
	SalesSummary   SalesSummary   `json:"sales_summary"`
	CostSummary    CostSummary    `json:"cost_summary"`
	ExpenseSummary ExpenseSummary `json:"expense_summary"`
	ProfitAnalysis ProfitAnalysis `json:"profit_analysis"`
	Details        *ReportDetails `json:"details,omitempty"` // 使用指针，omitempty表示可选
}

// SalesSummary 销售汇总
type SalesSummary struct {
	TotalAmount   float64 `json:"total_amount"`    // 总销售额
	NetAmount     float64 `json:"net_amount"`      // 净销售额（扣除退款）
	OrderCount    int     `json:"order_count"`     // 订单数
	ItemCount     int     `json:"item_count"`      // 商品件数
	AvgOrderValue float64 `json:"avg_order_value"` // 平均客单价
	RefundAmount  float64 `json:"refund_amount"`   // 退款金额
	RefundCount   int     `json:"refund_count"`    // 退款笔数
}

// CostSummary 成本汇总
type CostSummary struct {
	ProductCost  float64 `json:"product_cost"`  // 产品成本
	ShippingCost float64 `json:"shipping_cost"` // 运费成本
	HandlingCost float64 `json:"handling_cost"` // 操作费用
	OtherCost    float64 `json:"other_cost"`    // 其他费用
	TotalCost    float64 `json:"total_cost"`    // 总成本
}

// ExpenseSummary 费用汇总
type ExpenseSummary struct {
	PromotionExpenses float64 `json:"promotion_expenses"` // 推广费用
	//PlatformFees      float64 `json:"platform_fees"`      // 平台费用
	OtherExpenses float64 `json:"other_expenses"` // 其他费用
	TotalExpenses float64 `json:"total_expenses"` // 总费用
}

// ProfitAnalysis 利润分析
type ProfitAnalysis struct {
	GrossProfit     float64 `json:"gross_profit"`     // 毛利润
	OperatingProfit float64 `json:"operating_profit"` // 营业利润
	NetProfit       float64 `json:"net_profit"`       // 净利润
	GrossMargin     string  `json:"gross_margin"`     // 毛利率（百分比字符串）
	NetMargin       string  `json:"net_margin"`       // 净利率（百分比字符串）
}

// ReportDetails 报表明细数据
type ReportDetails struct {
	TopProducts       []TopProduct       `json:"top_products,omitempty"`
	PromotionChannels []PromotionChannel `json:"promotion_channels,omitempty"`
}

// TopProduct 热销商品
type TopProduct struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	SalesAmount float64 `json:"sales_amount"`
	Quantity    int     `json:"quantity"`
}

// PromotionChannel 推广渠道
type PromotionChannel struct {
	Channel string  `json:"channel"`
	Expense float64 `json:"expense"`
	Orders  int     `json:"orders"`
}

// API响应结构
type ReportResponse struct {
	Success   bool         `json:"success"`
	Data      *DailyReport `json:"data"`
	Timestamp string       `json:"timestamp"`
}

// ProductCostDetail 商品成本详情
type ProductCostDetail struct {
	PurchasePrice float64 `json:"purchase_price"` // 采购价格
	ShippingCost  float64 `json:"shipping_cost"`  // 运费
	HandlingCost  float64 `json:"handling_cost"`  // 操作费
	OtherCost     float64 `json:"other_cost"`     // 其他费用
	TotalCost     float64 `json:"total_cost"`     // 总成本
	Found         bool    `json:"found"`          // 是否找到匹配
}
