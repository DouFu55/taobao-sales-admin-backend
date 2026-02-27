package validators

type PromotionExpenseValidator struct {
	Id            int64   `json:"id,omitempty" form:"id"`                                              //
	PromotionDate string  `json:"promotion_date" form:"promotion_date" binding:"required"`             // 推广日期
	PromotionType int     `json:"promotion_type" form:"promotion_type" binding:"required,min=1,max=7"` // 推广类型：1-7
	PlanID        string  `json:"plan_id" form:"plan_id" binding:"max=100"`                            // 计划ID/名称
	ExpenseAmount float64 `json:"expense_amount" form:"expense_amount" binding:"required,min=0"`       // 支出金额(元)
	OrdersCount   int     `json:"orders_count" form:"orders_count" binding:"min=0"`                    // 成交笔数
	Remark        string  `json:"remark" form:"remark" binding:"max=200"`                              // 备注
}
