package validators

type ProductCostValidator struct {
	ID            *int64   `json:"id" form:"id" validate:"gt=0"`
	SKU           string   `json:"sku" form:"sku" validate:"required,max=100"`
	ProductID     int64    `json:"product_id" form:"product_id" validate:"required,gt=0"`
	ProductName   string   `json:"product_name" form:"product_name" validate:"required,max=200"`
	SpecInfo      string   `json:"spec_info" form:"spec_info" validate:"max=500"`
	PurchasePrice *float64 `json:"purchase_price" form:"purchase_price" validate:"required,gte=0"`
	ShippingCost  *float64 `json:"shipping_cost" form:"shipping_cost" validate:"required,gte=0"`
	HandlingCost  *float64 `json:"handling_cost" form:"handling_cost" validate:"required,gte=0"`
	OtherCost     *float64 `json:"other_cost" form:"other_cost" validate:"required,gte=0"`
	TotalCost     *float64 `json:"total_cost" form:"total_cost" validate:"required,gte=0"`
	Supplier      string   `json:"supplier" form:"supplier" validate:"required,max=200"`
	Stock         *int     `json:"stock" form:"stock" validate:"required,gte=0"`
	IsActive      *bool    `json:"is_active" form:"is_active" validate:"required"`
}
