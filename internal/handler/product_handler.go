package handler

import (
	"net/http"

	"github.com/imran4u/Oolio/internal/repository"

	"github.com/gin-gonic/gin"
)

func ListProducts(c *gin.Context) {

	products := repository.GetAllProducts()

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {

	id := c.Param("productId")

	product, found := repository.GetProductByID(id)

	if !found {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "product not found",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}
