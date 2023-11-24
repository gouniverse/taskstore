package taskstore

import (
	"time"
)

const (
	TaskStatusActive   = "active"
	TaskStatusCanceled = "canceled"
)

// Task type represents a definition of a task
type Task struct {
	ID          string     `json:"id" db:"id"`                   // varchar(40)  primary_key
	Status      string     `json:"status" db:"status"`           // varchar(40)  NOT NULL
	Alias       string     `json:"alias" db:"alias"`             // varchar(40)  NOT NULL
	Title       string     `json:"title" db:"title"`             // varchar(255) NOT NULL
	Description string     `json:"description" db:"description"` // text         DEFAULT NULL
	Memo        string     `json:"memo" db:"memo"`               // text         DEFAULT NULL
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`   // datetime     NOT NULL
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`   // datetime     NOT NULL
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`   // datetime     DEFAULT NULL
}
