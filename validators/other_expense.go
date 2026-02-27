package validators

// OtherExpenseValidator 其他支出验证器（简化版）
type OtherExpenseValidator struct {
	ID          int64   `json:"id,omitempty"`                                // 主键ID
	ExpenseDate string  `json:"expense_date" binding:"required"`             // 支出日期，格式：YYYY-MM-DD
	ExpenseType int     `json:"expense_type" binding:"required,min=1,max=7"` // 支出类型：1-7
	ItemName    string  `json:"item_name" binding:"required,max=200"`        // 项目名称
	Amount      float64 `json:"amount" binding:"required,min=0"`             // 金额(元)
	SubOrderNo  string  `json:"sub_order_no,omitempty" binding:"max=100"`    // 子订单号（选填）
	Remark      string  `json:"remark,omitempty" binding:"max=500"`          // 备注
}
