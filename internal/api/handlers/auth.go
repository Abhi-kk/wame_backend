package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"wati-clone-backend/internal/core/models"
	"wati-clone-backend/internal/infrastructure/db"
	"wati-clone-backend/internal/pkg/utils"
)

type SignupRequest struct {
	CompanyName string `json:"company_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context) {
	var req SignupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Create Tenant
	tenant := models.Tenant{Name: req.CompanyName}
	if err := db.DB.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tenant"})
		return
	}

	// 2. Hash Password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 3. Create User
	user := models.User{
		TenantID: tenant.ID,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "admin",
	}
	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// 4. Generate Token
	token, _ := utils.GenerateToken(user.ID, tenant.ID, user.Role)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Signup successful",
		"token":   token,
		"tenant":  tenant,
		"user":    user,
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(user.ID, user.TenantID, user.Role)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
