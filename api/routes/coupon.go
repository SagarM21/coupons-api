package routes

import (
	"net/http"

	"e-commerce/pkg/models"
	"e-commerce/pkg/services"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	Service *services.CouponService
}

func RegisterCouponRoutes(r *gin.Engine, handler *CouponHandler) {
	r.POST("/coupons", handler.CreateCoupon)
	r.GET("/coupons", handler.GetAllCoupons)
	r.GET("/coupons/:id", handler.GetCoupon)
	r.PUT("/coupons/:id", handler.UpdateCoupon)
	r.DELETE("/coupons/:id", handler.DeleteCoupon)
	r.POST("/applicable-coupons", handler.GetApplicableCoupons)
	r.POST("/apply-coupon/:id", handler.ApplyCoupon)
}

func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var input services.CouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	coupon, err := h.Service.Create(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, coupon)
}

func (h *CouponHandler) GetAllCoupons(c *gin.Context) {
	coupons, err := h.Service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch coupons"})
		return
	}

	if coupons == nil {
		coupons = []models.Coupon{}
	}

	c.JSON(http.StatusOK, coupons)
}

func (h *CouponHandler) GetCoupon(c *gin.Context) {
	id := c.Param("id")

	coupon, err := h.Service.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "coupon not found"})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

func (h *CouponHandler) UpdateCoupon(c *gin.Context) {
	id := c.Param("id")

	var input services.CouponInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	coupon, err := h.Service.Update(id, input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coupon)
}

func (h *CouponHandler) DeleteCoupon(c *gin.Context) {
	id := c.Param("id")

	if err := h.Service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "coupon deleted"})
}

func (h *CouponHandler) GetApplicableCoupons(c *gin.Context) {
	var req struct {
		Cart models.Cart `json:"cart"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.Cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart must have at least one item"})
		return
	}

	applicable, err := h.Service.GetApplicableCoupons(req.Cart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get applicable coupons"})
		return
	}

	if applicable == nil {
		applicable = []models.ApplicableCoupon{}
	}

	c.JSON(http.StatusOK, gin.H{"applicable_coupons": applicable})
}

func (h *CouponHandler) ApplyCoupon(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Cart models.Cart `json:"cart"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.Cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart must have at least one item"})
		return
	}

	updatedCart, err := h.Service.ApplyCoupon(id, req.Cart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"updated_cart": updatedCart})
}
