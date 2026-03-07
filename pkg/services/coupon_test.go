package services

import (
	"encoding/json"
	"testing"

	"e-commerce/pkg/models"
)

func TestCartWiseCalculate(t *testing.T) {
	s := &CartWiseStrategy{}

	tests := []struct {
		name      string
		details   models.CartWiseDetails
		cart      models.Cart
		wantDisc  float64
		wantApply bool
	}{
		{
			name:    "cart total exceeds threshold",
			details: models.CartWiseDetails{Threshold: 100, Discount: 10},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 3, Price: 50},
			}},
			wantDisc:  15,
			wantApply: true,
		},
		{
			name:    "cart total equals threshold - not applicable",
			details: models.CartWiseDetails{Threshold: 150, Discount: 10},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 3, Price: 50},
			}},
			wantDisc:  0,
			wantApply: false,
		},
		{
			name:    "cart total below threshold",
			details: models.CartWiseDetails{Threshold: 200, Discount: 10},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 1, Price: 50},
			}},
			wantDisc:  0,
			wantApply: false,
		},
		{
			name:    "multiple items exceeding threshold",
			details: models.CartWiseDetails{Threshold: 100, Discount: 10},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 6, Price: 50},
				{ProductID: 2, Quantity: 3, Price: 30},
				{ProductID: 3, Quantity: 2, Price: 25},
			}},
			wantDisc:  44,
			wantApply: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detailsJSON, _ := json.Marshal(tt.details)
			discount, applicable := s.Calculate(string(detailsJSON), tt.cart)
			if applicable != tt.wantApply {
				t.Errorf("applicable = %v, want %v", applicable, tt.wantApply)
			}
			if discount != tt.wantDisc {
				t.Errorf("discount = %v, want %v", discount, tt.wantDisc)
			}
		})
	}
}

func TestProductWiseCalculate(t *testing.T) {
	s := &ProductWiseStrategy{}

	tests := []struct {
		name      string
		details   models.ProductWiseDetails
		cart      models.Cart
		wantDisc  float64
		wantApply bool
	}{
		{
			name:    "product in cart",
			details: models.ProductWiseDetails{ProductID: 1, Discount: 20},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 2, Price: 100},
			}},
			wantDisc:  40,
			wantApply: true,
		},
		{
			name:    "product not in cart",
			details: models.ProductWiseDetails{ProductID: 5, Discount: 20},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 2, Price: 100},
			}},
			wantDisc:  0,
			wantApply: false,
		},
		{
			name:    "single quantity product",
			details: models.ProductWiseDetails{ProductID: 1, Discount: 50},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 1, Price: 200},
			}},
			wantDisc:  100,
			wantApply: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detailsJSON, _ := json.Marshal(tt.details)
			discount, applicable := s.Calculate(string(detailsJSON), tt.cart)
			if applicable != tt.wantApply {
				t.Errorf("applicable = %v, want %v", applicable, tt.wantApply)
			}
			if discount != tt.wantDisc {
				t.Errorf("discount = %v, want %v", discount, tt.wantDisc)
			}
		})
	}
}

func TestBxGyCalculate(t *testing.T) {
	s := &BxGyStrategy{}

	tests := []struct {
		name      string
		details   models.BxGyDetails
		cart      models.Cart
		wantDisc  float64
		wantApply bool
	}{
		{
			name: "basic b2g1 - applicable",
			details: models.BxGyDetails{
				BuyProducts:     []models.ProductQuantity{{ProductID: 1, Quantity: 2}},
				GetProducts:     []models.ProductQuantity{{ProductID: 3, Quantity: 1}},
				RepetitionLimit: 1,
			},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 3, Price: 50},
				{ProductID: 3, Quantity: 1, Price: 25},
			}},
			wantDisc:  25,
			wantApply: true,
		},
		{
			name: "b2g1 - not enough buy products",
			details: models.BxGyDetails{
				BuyProducts:     []models.ProductQuantity{{ProductID: 1, Quantity: 2}},
				GetProducts:     []models.ProductQuantity{{ProductID: 3, Quantity: 1}},
				RepetitionLimit: 1,
			},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 1, Price: 50},
				{ProductID: 3, Quantity: 1, Price: 25},
			}},
			wantDisc:  0,
			wantApply: false,
		},
		{
			name: "b2g1 with repetition limit",
			details: models.BxGyDetails{
				BuyProducts: []models.ProductQuantity{
					{ProductID: 1, Quantity: 3},
					{ProductID: 2, Quantity: 3},
				},
				GetProducts:     []models.ProductQuantity{{ProductID: 3, Quantity: 1}},
				RepetitionLimit: 2,
			},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 6, Price: 50},
				{ProductID: 2, Quantity: 3, Price: 30},
				{ProductID: 3, Quantity: 2, Price: 25},
			}},
			wantDisc:  25,
			wantApply: true,
		},
		{
			name: "multiple get products",
			details: models.BxGyDetails{
				BuyProducts:     []models.ProductQuantity{{ProductID: 1, Quantity: 2}},
				GetProducts:     []models.ProductQuantity{{ProductID: 2, Quantity: 1}, {ProductID: 3, Quantity: 1}},
				RepetitionLimit: 1,
			},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 2, Price: 50},
				{ProductID: 2, Quantity: 1, Price: 30},
				{ProductID: 3, Quantity: 1, Price: 25},
			}},
			wantDisc:  55,
			wantApply: true,
		},
		{
			name: "get product not in cart",
			details: models.BxGyDetails{
				BuyProducts:     []models.ProductQuantity{{ProductID: 1, Quantity: 2}},
				GetProducts:     []models.ProductQuantity{{ProductID: 5, Quantity: 1}},
				RepetitionLimit: 1,
			},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 2, Price: 50},
			}},
			wantDisc:  0,
			wantApply: false,
		},
		{
			name: "free qty capped by cart quantity",
			details: models.BxGyDetails{
				BuyProducts:     []models.ProductQuantity{{ProductID: 1, Quantity: 1}},
				GetProducts:     []models.ProductQuantity{{ProductID: 2, Quantity: 3}},
				RepetitionLimit: 5,
			},
			cart: models.Cart{Items: []models.CartItem{
				{ProductID: 1, Quantity: 10, Price: 50},
				{ProductID: 2, Quantity: 2, Price: 25},
			}},
			wantDisc:  50,
			wantApply: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detailsJSON, _ := json.Marshal(tt.details)
			discount, applicable := s.Calculate(string(detailsJSON), tt.cart)
			if applicable != tt.wantApply {
				t.Errorf("applicable = %v, want %v", applicable, tt.wantApply)
			}
			if discount != tt.wantDisc {
				t.Errorf("discount = %v, want %v", discount, tt.wantDisc)
			}
		})
	}
}

func TestValidateViaRegistry(t *testing.T) {
	tests := []struct {
		name       string
		couponType models.CouponType
		details    json.RawMessage
		wantErr    bool
	}{
		{
			name:       "valid cart-wise",
			couponType: models.CouponTypeCartWise,
			details:    json.RawMessage(`{"threshold": 100, "discount": 10}`),
			wantErr:    false,
		},
		{
			name:       "cart-wise negative threshold",
			couponType: models.CouponTypeCartWise,
			details:    json.RawMessage(`{"threshold": -10, "discount": 10}`),
			wantErr:    true,
		},
		{
			name:       "cart-wise discount over 100",
			couponType: models.CouponTypeCartWise,
			details:    json.RawMessage(`{"threshold": 100, "discount": 150}`),
			wantErr:    true,
		},
		{
			name:       "valid product-wise",
			couponType: models.CouponTypeProductWise,
			details:    json.RawMessage(`{"product_id": 1, "discount": 20}`),
			wantErr:    false,
		},
		{
			name:       "product-wise zero product_id",
			couponType: models.CouponTypeProductWise,
			details:    json.RawMessage(`{"product_id": 0, "discount": 20}`),
			wantErr:    true,
		},
		{
			name:       "valid bxgy",
			couponType: models.CouponTypeBxGy,
			details: json.RawMessage(`{
				"buy_products": [{"product_id": 1, "quantity": 3}],
				"get_products": [{"product_id": 2, "quantity": 1}],
				"repition_limit": 2
			}`),
			wantErr: false,
		},
		{
			name:       "bxgy empty buy products",
			couponType: models.CouponTypeBxGy,
			details: json.RawMessage(`{
				"buy_products": [],
				"get_products": [{"product_id": 2, "quantity": 1}],
				"repition_limit": 2
			}`),
			wantErr: true,
		},
		{
			name:       "unknown coupon type",
			couponType: "invalid-type",
			details:    json.RawMessage(`{}`),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strategy, err := GetStrategy(tt.couponType)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("GetStrategy() unexpected error: %v", err)
				}
				return
			}
			err = strategy.Validate(tt.details)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCartTotal(t *testing.T) {
	cart := models.Cart{Items: []models.CartItem{
		{ProductID: 1, Quantity: 6, Price: 50},
		{ProductID: 2, Quantity: 3, Price: 30},
		{ProductID: 3, Quantity: 2, Price: 25},
	}}

	total := cartTotal(cart)
	expected := 440.0
	if total != expected {
		t.Errorf("cartTotal() = %v, want %v", total, expected)
	}
}
