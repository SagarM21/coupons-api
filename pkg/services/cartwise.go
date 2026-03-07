package services

import (
	"encoding/json"
	"fmt"
	"math"

	"e-commerce/pkg/models"
)

func init() {
	RegisterStrategy(models.CouponTypeCartWise, &CartWiseStrategy{})
}

type CartWiseStrategy struct{}

func (s *CartWiseStrategy) Validate(details json.RawMessage) error {
	var d models.CartWiseDetails
	if err := json.Unmarshal(details, &d); err != nil {
		return fmt.Errorf("invalid cart-wise details: %w", err)
	}
	if d.Threshold < 0 {
		return fmt.Errorf("threshold must be non-negative")
	}
	if d.Discount <= 0 || d.Discount > 100 {
		return fmt.Errorf("discount must be between 0 and 100")
	}
	return nil
}

func (s *CartWiseStrategy) Calculate(details string, cart models.Cart) (float64, bool) {
	var d models.CartWiseDetails
	if err := json.Unmarshal([]byte(details), &d); err != nil {
		return 0, false
	}

	total := cartTotal(cart)
	if total <= d.Threshold {
		return 0, false
	}

	discount := math.Round(total*d.Discount) / 100
	return discount, true
}

func (s *CartWiseStrategy) Apply(details string, cart models.Cart) ([]models.CartItem, float64) {
	var d models.CartWiseDetails
	json.Unmarshal([]byte(details), &d)

	totalPrice := cartTotal(cart)
	totalDiscount := math.Round(totalPrice*d.Discount) / 100

	items := make([]models.CartItem, len(cart.Items))
	copy(items, cart.Items)

	for i := range items {
		itemTotal := float64(items[i].Quantity) * items[i].Price
		items[i].TotalDiscount = math.Round(totalDiscount*itemTotal/totalPrice*100) / 100
	}

	return items, totalDiscount
}
