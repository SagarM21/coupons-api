package services

import (
	"encoding/json"
	"fmt"
	"math"

	"e-commerce/pkg/models"
)

func init() {
	RegisterStrategy(models.CouponTypeProductWise, &ProductWiseStrategy{})
}

type ProductWiseStrategy struct{}

func (s *ProductWiseStrategy) Validate(details json.RawMessage) error {
	var d models.ProductWiseDetails
	if err := json.Unmarshal(details, &d); err != nil {
		return fmt.Errorf("invalid product-wise details: %w", err)
	}
	if d.ProductID <= 0 {
		return fmt.Errorf("product_id must be positive")
	}
	if d.Discount <= 0 || d.Discount > 100 {
		return fmt.Errorf("discount must be between 0 and 100")
	}
	return nil
}

func (s *ProductWiseStrategy) Calculate(details string, cart models.Cart) (float64, bool) {
	var d models.ProductWiseDetails
	if err := json.Unmarshal([]byte(details), &d); err != nil {
		return 0, false
	}

	for _, item := range cart.Items {
		if item.ProductID == d.ProductID {
			discount := math.Round(float64(item.Quantity)*item.Price*d.Discount) / 100
			return discount, true
		}
	}
	return 0, false
}

func (s *ProductWiseStrategy) Apply(details string, cart models.Cart) ([]models.CartItem, float64) {
	var d models.ProductWiseDetails
	json.Unmarshal([]byte(details), &d)

	items := make([]models.CartItem, len(cart.Items))
	copy(items, cart.Items)

	var totalDiscount float64
	for i := range items {
		if items[i].ProductID == d.ProductID {
			items[i].TotalDiscount = math.Round(float64(items[i].Quantity)*items[i].Price*d.Discount) / 100
			totalDiscount += items[i].TotalDiscount
		}
	}

	return items, totalDiscount
}
