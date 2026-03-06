package handler

import (
	"net/http"

	"github.com/imran4u/Oolio/internal/model"
	"github.com/imran4u/Oolio/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	CouponService interface {
		ValidateCoupon(string) bool
	}
}

func NewOrderHandler(couponSvc interface {
	ValidateCoupon(string) bool
}) *OrderHandler {
	return &OrderHandler{couponSvc}
}

func (h *OrderHandler) PlaceOrder(c *gin.Context) {

	var req model.OrderReq

	if err := c.ShouldBindJSON(&req); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if req.CouponCode != "" {

		valid := h.CouponService.ValidateCoupon(req.CouponCode)

		if !valid {
			c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": "invalid coupon",
			})
			return
		}
	}

	var products []model.Product

	for _, item := range req.Items {

		p, found := repository.GetProductByID(item.ProductID)

		if !found {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "invalid product",
			})
			return
		}

		products = append(products, *p)
	}

	order := model.Order{
		ID:       uuid.New().String(),
		Items:    req.Items,
		Products: products,
	}

	c.JSON(http.StatusOK, order)
}
