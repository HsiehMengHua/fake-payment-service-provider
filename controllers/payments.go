package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PaymentController interface {
	PayInPage(c *gin.Context)
	PayOutPage(c *gin.Context)
	ConfirmPayment(c *gin.Context)
	CancelPayment(c *gin.Context)
}

type TransactionData struct {
	TransactionID   string
	ConfirmCallback string
	CancelCallback  string
}

type paymentController struct {
}

type callbackRequest struct {
	TransactionID string `json:"transaction_id"`
	CallbackURL   string `json:"callback_url"`
}

func NewPaymentController() PaymentController {
	return &paymentController{}
}

func (ctrl *paymentController) PayInPage(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	amount := c.Query("amount")
	merchant := c.Query("merchant")
	confirmCallback := c.Query("confirm_callback")
	cancelCallback := c.Query("cancel_callback")

	c.HTML(http.StatusOK, "payment.html", gin.H{
		"TransactionID":   transactionID,
		"Amount":          amount,
		"Merchant":        merchant,
		"ConfirmCallback": confirmCallback,
		"CancelCallback":  cancelCallback,
	})
}

func (ctrl *paymentController) PayOutPage(c *gin.Context) {
	transactionID := c.Param("transaction_id")
	amount := c.Query("amount")
	bankCode := c.Query("bank_code")
	accountNumber := c.Query("account_number")
	confirmCallback := c.Query("confirm_callback")
	cancelCallback := c.Query("cancel_callback")

	c.HTML(http.StatusOK, "payout.html", gin.H{
		"TransactionID":   transactionID,
		"Amount":          amount,
		"BankCode":        bankCode,
		"AccountNumber":   accountNumber,
		"ConfirmCallback": confirmCallback,
		"CancelCallback":  cancelCallback,
	})
}

func (ctrl *paymentController) ConfirmPayment(c *gin.Context) {
	var request callbackRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sendCallbackRequest(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment confirmed"})
}

func (ctrl *paymentController) CancelPayment(c *gin.Context) {
	var request callbackRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sendCallbackRequest(request); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment canceled"})
}

func sendCallbackRequest(model callbackRequest) error {
	payload := map[string]string{
		"transaction_id": model.TransactionID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		log.Panic("Failed to marshal callback payload: ", err)
	}
	req, err := http.NewRequest("POST", model.CallbackURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Panic("Failed to create callback request: ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", os.Getenv("PSP_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Panic("Failed to send callback request: ", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		err := fmt.Errorf("callback request failed. Response status: %s, message: %s", resp.Status, body)
		log.Error(err.Error())
		return err
	}

	return nil
}
