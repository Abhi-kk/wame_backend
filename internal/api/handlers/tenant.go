package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wati-clone-backend/internal/core/models"
	"wati-clone-backend/internal/infrastructure/db"
)

type CreateTenantRequest struct {
	Name string `json:"name" binding:"required"`
}

func CreateTenant(c *gin.Context) {
	var req CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant := models.Tenant{Name: req.Name}
	if err := db.DB.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
		return
	}

	c.JSON(http.StatusCreated, tenant)
}
