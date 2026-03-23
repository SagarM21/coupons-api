package models

import "time"

type CouponType string

const (
	CouponTypeCartWise    CouponType = "cart-wise"
	CouponTypeProductWise CouponType = "product-wise"
	CouponTypeBxGy        CouponType = "bxgy"
)


type Coupon struct {
	ID        string     `json:"id" cql:"id"`
	Type      CouponType `json:"type" cql:"type"`
	Details   string     `json:"details" cql:"details"` // JSON string stored in DB
	CreatedAt time.Time  `json:"created_at" cql:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" cql:"expires_at"`
}

// CartWiseDetails holds details for cart-wise coupons.
type CartWiseDetails struct {
	Threshold float64 `json:"threshold"`
	Discount  float64 `json:"discount"` // percentage
}

// ProductWiseDetails holds details for product-wise coupons.
type ProductWiseDetails struct {
	ProductID int     `json:"product_id"`
	Discount  float64 `json:"discount"` // percentage
}

// BxGyDetails holds details for BxGy coupons.
type BxGyDetails struct {
	BuyProducts    []ProductQuantity `json:"buy_products"`
	GetProducts    []ProductQuantity `json:"get_products"`
	RepetitionLimit int              `json:"repition_limit"`
}

type ProductQuantity struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
