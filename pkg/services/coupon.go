package services

import (
	"encoding/json"
	"fmt"
	"time"

	"e-commerce/pkg/db"
	"e-commerce/pkg/models"

	"github.com/google/uuid"
)

type CouponService struct {
	DB *db.DB
}

type CouponInput struct {
	Type      models.CouponType `json:"type"`
	Details   json.RawMessage   `json:"details"`
	ExpiresAt *time.Time        `json:"expires_at,omitempty"`
}

func (s *CouponService) Create(input CouponInput) (*models.Coupon, error) {
	strategy, err := GetStrategy(input.Type)
	if err != nil {
		return nil, err
	}

	if err := strategy.Validate(input.Details); err != nil {
		return nil, err
	}

	coupon := models.Coupon{
		ID:        uuid.New().String(),
		Type:      input.Type,
		Details:   string(input.Details),
		CreatedAt: time.Now(),
		ExpiresAt: input.ExpiresAt,
	}

	if err := s.DB.CreateCoupon(coupon); err != nil {
		return nil, fmt.Errorf("failed to create coupon: %w", err)
	}
	return &coupon, nil
}

func (s *CouponService) Get(id string) (*models.Coupon, error) {
	return s.DB.GetCoupon(id)
}

func (s *CouponService) GetAll() ([]models.Coupon, error) {
	return s.DB.GetAllCoupons()
}

func (s *CouponService) Update(id string, input CouponInput) (*models.Coupon, error) {
	existing, err := s.DB.GetCoupon(id)
	if err != nil {
		return nil, err
	}

	strategy, err := GetStrategy(input.Type)
	if err != nil {
		return nil, err
	}

	if err := strategy.Validate(input.Details); err != nil {
		return nil, err
	}

	existing.Type = input.Type
	existing.Details = string(input.Details)
	existing.ExpiresAt = input.ExpiresAt

	if err := s.DB.UpdateCoupon(*existing); err != nil {
		return nil, fmt.Errorf("failed to update coupon: %w", err)
	}
	return existing, nil
}

func (s *CouponService) Delete(id string) error {
	return s.DB.DeleteCoupon(id)
}

func (s *CouponService) GetApplicableCoupons(cart models.Cart) ([]models.ApplicableCoupon, error) {
	coupons, err := s.DB.GetAllCoupons()
	if err != nil {
		return nil, err
	}

	var applicable []models.ApplicableCoupon

	for _, c := range coupons {
		if s.DB.IsExpired(&c) {
			continue
		}

		strategy, err := GetStrategy(c.Type)
		if err != nil {
			continue
		}

		discount, ok := strategy.Calculate(c.Details, cart)
		if ok && discount > 0 {
			applicable = append(applicable, models.ApplicableCoupon{
				CouponID: c.ID,
				Type:     c.Type,
				Discount: discount,
			})
		}
	}

	return applicable, nil
}

func (s *CouponService) ApplyCoupon(couponID string, cart models.Cart) (*models.UpdatedCart, error) {
	coupon, err := s.DB.GetCoupon(couponID)
	if err != nil {
		return nil, fmt.Errorf("coupon not found")
	}

	if s.DB.IsExpired(coupon) {
		return nil, fmt.Errorf("coupon has expired")
	}

	strategy, err := GetStrategy(coupon.Type)
	if err != nil {
		return nil, err
	}

	discount, ok := strategy.Calculate(coupon.Details, cart)
	if !ok || discount <= 0 {
		return nil, fmt.Errorf("coupon is not applicable to this cart")
	}

	items, totalDiscount := strategy.Apply(coupon.Details, cart)
	totalPrice := cartTotal(cart)

	return &models.UpdatedCart{
		Items:         items,
		TotalPrice:    totalPrice,
		TotalDiscount: totalDiscount,
		FinalPrice:    totalPrice - totalDiscount,
	}, nil
}

func cartTotal(cart models.Cart) float64 {
	var total float64
	for _, item := range cart.Items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
