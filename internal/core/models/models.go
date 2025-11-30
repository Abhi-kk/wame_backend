package models

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// Settings (JSONB in real app, simple string for MVP)
	WhatsAppToken    string `json:"-"`
	WhatsAppPhoneID  string `json:"-"`
}

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"index" json:"tenant_id"`
	Email     string         `gorm:"uniqueIndex" json:"email"`
	Password  string         `json:"-"`
	Role      string         `json:"role"` // admin, agent
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Contact struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"index" json:"tenant_id"`
	Name      string         `json:"name"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Message struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	TenantID       uint      `gorm:"index" json:"tenant_id"`
	ContactID      uint      `gorm:"index" json:"contact_id"`
	Channel        string    `json:"channel"` // whatsapp, email, sms
	Direction      string    `json:"direction"` // inbound, outbound
	Content        string    `json:"content"`
	Status         string    `json:"status"` // sent, delivered, read, failed
	ExternalID     string    `json:"external_id"` // ID from provider
	CreatedAt      time.Time `json:"created_at"`
}
