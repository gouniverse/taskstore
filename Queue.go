package taskstore

import (
	"encoding/json"
	"time"

	"github.com/golang-module/carbon/v2"
)

const (
	QueueStatusCanceled = "canceled"
	QueueStatusDeleted  = "deleted"
	QueueStatusFailed   = "failed"
	QueueStatusPaused   = "paused"
	QueueStatusQueued   = "queued"
	QueueStatusRunning  = "running"
	QueueStatusSuccess  = "success"
)

// Queue type represents an queued task in the queue
type Queue struct {
	ID          string     `json:"id" db:"id"`                     // varchar (40) primary_key
	Status      string     `json:"status" db:"status"`             // varchar(40) DEFAULT 'queued'
	TaskID      string     `json:"task_id" db:"task_id"`           // varchar(40)
	Parameters  string     `json:"parameters" db:"parameters"`     // text
	Output      string     `json:"output" db:"output"`             // text
	Details     string     `json:"details" db:"details"`           // text
	Attempts    int        `json:"attempts" db:"attempts"`         // int
	StartedAt   *time.Time `json:"started_at" db:"started_at"`     // datetime DEFAULT NULL
	CompletedAt *time.Time `json:"completed_at" db:"completed_at"` // datetime DEFAULT NULL
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`     // datetime NOT NULL
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`     // datetime NOT NULL
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`     // datetime DEFAULT NULL
}

// TableName the name of the queue table
// func (Queue) TableName() string {
// 	return "snv_tasks_queue"
// }

// AppendDetails appends details to the queued task
// !!! warning does not auto-save it for performance reasons
func (queuedTask *Queue) AppendDetails(details string) {
	ts := carbon.Now().Format("Y-m-d H:i:s")
	text := queuedTask.Details
	if text != "" {
		text += "\n"
	}
	text += ts + " : " + details
	queuedTask.Details = text
}

// GetParameters gets the parameters of the queued task
func (queuedQueue *Queue) GetParameters() (map[string]interface{}, error) {
	var parameters map[string]interface{}
	jsonErr := json.Unmarshal([]byte(queuedQueue.Parameters), &parameters)
	if jsonErr != nil {
		return parameters, jsonErr
	}
	return parameters, nil
}
