package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PaymentController interface {
	PaymentPage(c *gin.Context)
}

type paymentController struct {
}

func NewPaymentController() PaymentController {
	return &paymentController{}
}

func (ctrl *paymentController) PaymentPage(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	amount := c.Param("amount")
	merchant := c.Param("merchant")

	c.HTML(http.StatusOK, "payment.html", gin.H{
		"TransactionID": transactionID,
		"Amount":        amount,
		"Merchant":      merchant,
	})
}
