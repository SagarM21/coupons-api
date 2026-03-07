package db

import (
	"e-commerce/pkg/models"
	"time"
)

func (d *DB) CreateCoupon(c models.Coupon) error {
	return d.Session.Query(
		`INSERT INTO coupons (id, type, details, created_at, expires_at) VALUES (?, ?, ?, ?, ?)`,
		c.ID, string(c.Type), c.Details, c.CreatedAt, c.ExpiresAt,
	).Exec()
}

func (d *DB) GetCoupon(id string) (*models.Coupon, error) {
	var c models.Coupon
	var couponType string
	err := d.Session.Query(
		`SELECT id, type, details, created_at, expires_at FROM coupons WHERE id = ?`, id,
	).Scan(&c.ID, &couponType, &c.Details, &c.CreatedAt, &c.ExpiresAt)
	if err != nil {
		return nil, err
	}
	c.Type = models.CouponType(couponType)
	return &c, nil
}

func (d *DB) GetAllCoupons() ([]models.Coupon, error) {
	var coupons []models.Coupon
	iter := d.Session.Query(`SELECT id, type, details, created_at, expires_at FROM coupons`).Iter()

	var c models.Coupon
	var couponType string
	for iter.Scan(&c.ID, &couponType, &c.Details, &c.CreatedAt, &c.ExpiresAt) {
		c.Type = models.CouponType(couponType)
		coupons = append(coupons, c)
		c = models.Coupon{}
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return coupons, nil
}

func (d *DB) UpdateCoupon(c models.Coupon) error {
	return d.Session.Query(
		`UPDATE coupons SET type = ?, details = ?, expires_at = ? WHERE id = ?`,
		string(c.Type), c.Details, c.ExpiresAt, c.ID,
	).Exec()
}

func (d *DB) DeleteCoupon(id string) error {
	return d.Session.Query(`DELETE FROM coupons WHERE id = ?`, id).Exec()
}

func (d *DB) IsExpired(c *models.Coupon) bool {
	if c.ExpiresAt == nil {
		return false
	}
	return c.ExpiresAt.Before(time.Now())
}
