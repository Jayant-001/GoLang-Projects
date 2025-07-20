package model

import "time"

type Task struct {
	ID          string        `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Status      string        `json:"status"`       // e.g., "pending", "in-progress", "completed"
	Delayed     time.Duration `json:"delayed_time"` // Unix timestamp for delayed execution
	CreatedAt   time.Time     `json:"created_at"`   // Unix timestamp for task creation
}

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}
