package router

import (
	"github.com/imran4u/Oolio/internal/handler"
	"github.com/imran4u/Oolio/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, couponSvc interface {
	ValidateCoupon(string) bool
}) {

	api := r.Group("/api")

	api.GET("/product", handler.ListProducts)
	api.GET("/product/:productId", handler.GetProduct)

	orderHandler := handler.NewOrderHandler(couponSvc)

	api.POST("/order",
		middleware.APIKeyAuth(),
		orderHandler.PlaceOrder,
	)
}
