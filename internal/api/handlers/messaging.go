package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wati-clone-backend/internal/core/models"
	"wati-clone-backend/internal/infrastructure/db"
)

type SendMessageRequest struct {
	ContactID uint   `json:"contact_id" binding:"required"`
	Channel   string `json:"channel" binding:"required"` // whatsapp, email, sms
	Content   string `json:"content" binding:"required"`
}

func SendMessage(c *gin.Context) {
	tenantID := c.GetUint("tenantID")
	var req SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Integrate with actual providers (WhatsApp, SendGrid, Twilio)
	// For MVP, just log to DB

	msg := models.Message{
		TenantID:  tenantID,
		ContactID: req.ContactID,
		Channel:   req.Channel,
		Direction: "outbound",
		Content:   req.Content,
		Status:    "sent", // Mock success
	}

	if err := db.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save message"})
		return
	}

	c.JSON(http.StatusOK, msg)
}

func GetMessages(c *gin.Context) {
	tenantID := c.GetUint("tenantID")
	var messages []models.Message
	
	// Simple pagination or limit could be added
	if err := db.DB.Where("tenant_id = ?", tenantID).Order("created_at desc").Limit(50).Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// WhatsAppWebhook handles incoming messages from Meta
func WhatsAppWebhook(c *gin.Context) {
	// TODO: Parse WhatsApp payload and save inbound message
	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func WhatsAppWebhookVerify(c *gin.Context) {
	// Verify token logic for Meta
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")

	if mode == "subscribe" && token == "mytesttoken" { // TODO: Env var
		c.String(http.StatusOK, challenge)
	} else {
		c.Status(http.StatusForbidden)
	}
}
