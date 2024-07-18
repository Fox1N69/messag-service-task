package models

import (
	"time"
)

// Client represents a client entity in the system.
type Client struct {
	ID          int64     `json:"id"`
	ClientName  string    `json:"client_name"`
	Version     int       `json:"version"`
	Image       string    `json:"image"`
	CPU         string    `json:"cpu"`
	Memory      string    `json:"memory"`
	Priority    float64   `json:"priority"`
	NeedRestart bool      `json:"need_restart"`
	SpawnedAt   time.Time `json:"spawned_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
