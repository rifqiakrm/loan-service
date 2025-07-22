package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRouter initializes all HTTP routes.
func SetupRouter(handler *Handler) *gin.Engine {
	r := gin.Default()

	r.GET("/loans", handler.ListLoans)
	r.GET("/loans/:id", handler.GetLoan)
	r.POST("/loans", handler.CreateLoan)
	r.POST("/loans/:id/approve", handler.ApproveLoan)
	r.POST("/loans/:id/invest", handler.InvestLoan)
	r.POST("/loans/:id/disburse", handler.DisburseLoan)

	return r
}
