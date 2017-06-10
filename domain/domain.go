package domain

import (
	"time"
)

type (
	Model struct {
		ID        uint      `json:"id" db:"id,omitempty" gorm:"primary_key"`
		Deleted   bool      `json:"deleted"`
		CreatedAt time.Time `json:"createdAt" db:"created_at"`
		UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
	}
)
