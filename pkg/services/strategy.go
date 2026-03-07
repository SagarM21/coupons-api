package services

import (
	"encoding/json"
	"fmt"

	"e-commerce/pkg/models"
)

// CouponStrategy defines the interface every coupon type must implement.
// To add a new coupon type: create a file, implement this interface, and register it via init().
type CouponStrategy interface {
	Validate(details json.RawMessage) error
	Calculate(details string, cart models.Cart) (float64, bool)
	Apply(details string, cart models.Cart) (items []models.CartItem, totalDiscount float64)
}

var registry = map[models.CouponType]CouponStrategy{}

func RegisterStrategy(couponType models.CouponType, strategy CouponStrategy) {
	registry[couponType] = strategy
}

func GetStrategy(couponType models.CouponType) (CouponStrategy, error) {
	s, ok := registry[couponType]
	if !ok {
		return nil, fmt.Errorf("unknown coupon type: %s", couponType)
	}
	return s, nil
}
