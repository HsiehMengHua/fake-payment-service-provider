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
	r.GET("/payment/:transaction_id/:amount/:merchant", ctrl.PaymentPage)

	return r
}
