package main

import (
	"github.com/imran4u/Oolio/internal/coupon"
	"github.com/imran4u/Oolio/internal/router"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// load coupons
	couponService := coupon.NewCouponService(
		"data/couponbase1.gz",
		"data/couponbase2.gz",
		"data/couponbase3.gz",
	)

	router.RegisterRoutes(r, couponService)

	r.Run(":8080")
}
