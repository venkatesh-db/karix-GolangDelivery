package models

import "time"

type UserGorm struct {
	ID        int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string      `gorm:"size:100" json:"name"`
	Email     string      `gorm:"size:100;uniqueIndex" json:"email"`
	Active    bool        `gorm:"default:true" json:"active"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	Orders    []OrderGorm `gorm:"foreignKey:UserID;references:ID" json:"orders"`
}

type OrderGorm struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"index" json:"user_id"`
	Amount    float64   `gorm:"type:decimal(10,2)" json:"amount"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

