package models

import "time"

const (
	StatusOpen       = "open"
	StatusInProgress = "in_progress"
	StatusClosed     = "closed"
)

type Ticket struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"not null;default:open"`
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	User        User      `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
