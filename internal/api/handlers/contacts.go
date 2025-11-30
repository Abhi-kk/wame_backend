package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wati-clone-backend/internal/core/models"
	"wati-clone-backend/internal/infrastructure/db"
)

type CreateContactRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func CreateContact(c *gin.Context) {
	tenantID := c.GetUint("tenantID")
	var req CreateContactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact := models.Contact{
		TenantID: tenantID,
		Name:     req.Name,
		Phone:    req.Phone,
		Email:    req.Email,
	}

	if err := db.DB.Create(&contact).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusCreated, contact)
}

func GetContacts(c *gin.Context) {
	tenantID := c.GetUint("tenantID")
	var contacts []models.Contact

	if err := db.DB.Where("tenant_id = ?", tenantID).Find(&contacts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}
