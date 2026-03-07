package models

type CartItem struct {
	ProductID     int     `json:"product_id"`
	Quantity      int     `json:"quantity"`
	Price         float64 `json:"price"`
	TotalDiscount float64 `json:"total_discount,omitempty"`
}

type Cart struct {
	Items []CartItem `json:"items"`
}

type ApplicableCoupon struct {
	CouponID string     `json:"coupon_id"`
	Type     CouponType `json:"type"`
	Discount float64    `json:"discount"`
}

type UpdatedCart struct {
	Items         []CartItem `json:"items"`
	TotalPrice    float64    `json:"total_price"`
	TotalDiscount float64    `json:"total_discount"`
	FinalPrice    float64    `json:"final_price"`
}
