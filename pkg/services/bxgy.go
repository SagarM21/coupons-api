package services

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"e-commerce/pkg/models"
)

func init() {
	RegisterStrategy(models.CouponTypeBxGy, &BxGyStrategy{})
}

type BxGyStrategy struct{}

func (s *BxGyStrategy) Validate(details json.RawMessage) error {
	var d models.BxGyDetails
	if err := json.Unmarshal(details, &d); err != nil {
		return fmt.Errorf("invalid bxgy details: %w", err)
	}
	if len(d.BuyProducts) == 0 {
		return fmt.Errorf("buy_products must not be empty")
	}
	if len(d.GetProducts) == 0 {
		return fmt.Errorf("get_products must not be empty")
	}
	if d.RepetitionLimit <= 0 {
		return fmt.Errorf("repition_limit must be positive")
	}
	return nil
}

func (s *BxGyStrategy) Calculate(details string, cart models.Cart) (float64, bool) {
	var d models.BxGyDetails
	if err := json.Unmarshal([]byte(details), &d); err != nil {
		return 0, false
	}

	cartMap := make(map[int]models.CartItem)
	for _, item := range cart.Items {
		cartMap[item.ProductID] = item
	}

	maxBuyReps := math.MaxInt32
	for _, bp := range d.BuyProducts {
		cartItem, exists := cartMap[bp.ProductID]
		if !exists {
			return 0, false
		}
		reps := cartItem.Quantity / bp.Quantity
		if reps < maxBuyReps {
			maxBuyReps = reps
		}
	}

	if maxBuyReps <= 0 {
		return 0, false
	}

	if maxBuyReps > d.RepetitionLimit {
		maxBuyReps = d.RepetitionLimit
	}

	var totalDiscount float64
	for _, gp := range d.GetProducts {
		cartItem, exists := cartMap[gp.ProductID]
		if !exists {
			continue
		}
		freeQty := gp.Quantity * maxBuyReps
		if freeQty > cartItem.Quantity {
			freeQty = cartItem.Quantity
		}
		totalDiscount += float64(freeQty) * cartItem.Price
	}

	return totalDiscount, totalDiscount > 0
}

func (s *BxGyStrategy) Apply(details string, cart models.Cart) ([]models.CartItem, float64) {
	var d models.BxGyDetails
	json.Unmarshal([]byte(details), &d)

	items := make([]models.CartItem, len(cart.Items))
	copy(items, cart.Items)

	cartMap := make(map[int]int) // productID -> index
	for i, item := range items {
		cartMap[item.ProductID] = i
	}

	maxBuyReps := math.MaxInt32
	for _, bp := range d.BuyProducts {
		idx, exists := cartMap[bp.ProductID]
		if !exists {
			return items, 0
		}
		reps := items[idx].Quantity / bp.Quantity
		if reps < maxBuyReps {
			maxBuyReps = reps
		}
	}
	if maxBuyReps > d.RepetitionLimit {
		maxBuyReps = d.RepetitionLimit
	}

	type getItem struct {
		idx     int
		price   float64
		freeQty int
	}
	var getItems []getItem
	for _, gp := range d.GetProducts {
		idx, exists := cartMap[gp.ProductID]
		if !exists {
			continue
		}
		freeQty := gp.Quantity * maxBuyReps
		if freeQty > items[idx].Quantity {
			freeQty = items[idx].Quantity
		}
		getItems = append(getItems, getItem{idx: idx, price: items[idx].Price, freeQty: freeQty})
	}

	sort.Slice(getItems, func(i, j int) bool {
		return getItems[i].price > getItems[j].price
	})

	var totalDiscount float64
	for _, gi := range getItems {
		discount := float64(gi.freeQty) * gi.price
		items[gi.idx].TotalDiscount = discount
		totalDiscount += discount
	}

	return items, totalDiscount
}
