package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/isirfanm/online-store/errorsx"
	"github.com/isirfanm/online-store/inventory"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/products/:sku", handleGetProduct)
	router.POST("/orders", handleCreateOrder)

	return router
}

func handleGetProduct(ctx *gin.Context) {
	SKU := ctx.Param("sku")

	p, err := inventory.Repo.FindProduct(SKU)
	if err != nil {
		errorsx.ErrorResponse(
			ctx,
			errorsx.NewGeneralErrorMsg(fmt.Sprintf("failed to get product %s. %s", SKU, err.Error())))
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func handleCreateOrder(ctx *gin.Context) {
	o := &inventory.OrderCreate{}
	if err := ctx.ShouldBindJSON(o); err != nil {
		errorsx.ErrorResponse(
			ctx,
			errorsx.NewGeneralErrorMsg("created order failed to parse request. "+err.Error()))
		return
	}

	ox, err := inventory.ProcessOrderCreate(o)
	if err != nil {
		errorsx.ErrorResponse(
			ctx,
			errorsx.NewGeneralErrorMsg("created order failed to process order. "+err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, ox)
}
