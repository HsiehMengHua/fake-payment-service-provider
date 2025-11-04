package router

import (
	"fake-payment-service-provider/controllers"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.StaticFile("/version", "./version.txt")

	ctrl := controllers.NewPaymentController()
	r.GET("/payin/:transaction_id", ctrl.PaymentPage)
	r.POST("/transaction/confirm", ctrl.ConfirmPayment)
	r.POST("/transaction/cancel", ctrl.CancelPayment)

	return r
}
