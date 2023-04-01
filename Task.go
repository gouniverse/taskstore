package taskstore

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/golang-module/carbon/v2"
	"github.com/gouniverse/uid"
)

const (
	TaskStatusCanceled = "canceled"
	TaskStatusDeleted  = "deleted"
	TaskStatusFailed   = "failed"
	TaskStatusPaused   = "paused"
	TaskStatusQueued   = "queued"
	TaskStatusRunning  = "running"
	TaskStatusSuccess  = "success"
)

// Task type
type Task struct {
	ID           string     `json:"id" db:"id"`                       // varchar (40) primary_key
	Status       string     `json:"status" db:"status"`               // varchar(40) DEFAULT 'queued'
	DefinitionID string     `json:"definition_id" db:"definition_id"` // varchar(40)
	Parameters   string     `json:"parameters" db:"parameters"`       // text
	Output       string     `json:"output" db:"output"`               // text
	Details      string     `json:"details" db:"details"`             // text
	Attempts     int        `json:"attempts" db:"attempts"`           // int
	StartedAt    *time.Time `json:"started_at" db:"started_at"`       // datetime DEFAULT NULL
	CompletedAt  *time.Time `json:"completed_at" db:"completed_at"`   // datetime DEFAULT NULL
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`       // datetime NOT NULL
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`       // datetime NOT NULL
	DeletedAt    *time.Time `json:"deleted_at" db:"deleted_at"`       // datetime DEFAULT NULL
}

// TableName the name of the Task table
func (Task) TableName() string {
	return "snv_tasks_task"
}

// AppendDetails appends details to the task, warning does not save it
func (queuedTask *Task) AppendDetails(details string) {
	ts := carbon.Now().Format("Y-m-d H:i:s")
	text := queuedTask.Details
	if text != "" {
		text += "\n"
	}
	text += ts + " : " + details
	queuedTask.Details = text
}

// GetParameters appends details to the task, warning does not save it
func (queuedTask *Task) GetParameters() (map[string]interface{}, error) {
	var parameters map[string]interface{}
	jsonErr := json.Unmarshal([]byte(queuedTask.Parameters), &parameters)
	if jsonErr != nil {
		return parameters, jsonErr
	}
	return parameters, nil
}

// TaskFail fails a task
func (st *Store) TaskFail(task *Task) error {
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	task.Status = TaskStatusFailed
	return st.TaskUpdate(task)
}

// TaskSuccess completes a task successfully
func (st *Store) TaskSuccess(task *Task) error {
	completedAt := time.Now()
	task.CompletedAt = &completedAt
	task.Status = TaskStatusSuccess
	return st.TaskUpdate(task)
}

// TaskCreate creates a Task
func (st *Store) TaskCreate(task *Task) error {
	if task.ID == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		task.ID = uid.MicroUid()
	}
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.taskTaskTableName).Rows(task).ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

type TaskListOptions struct {
	Status    string
	Limit     int
	SortBy    string
	SortOrder string
}

func (st *Store) TaskList(options TaskListOptions) ([]Task, error) {
	q := goqu.Dialect(st.dbDriverName).From(st.taskTaskTableName)

	if options.Status != "" {
		q = q.Where(goqu.C("status").Eq(options.Status))
	}

	if options.SortBy != "" {
		if options.SortOrder == "asc" {
			q = q.Order(goqu.I(options.SortBy).Asc())
		} else {
			q = q.Order(goqu.I(options.SortBy).Desc())
		}
	}

	q = q.Where(goqu.C("deleted_at").IsNull())

	if options.Limit != 0 {
		q = q.Limit(uint(options.Limit))
	}

	sqlStr, _, _ := q.Select().ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	list := []Task{}
	err := sqlscan.Select(context.Background(), st.db, &list, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		log.Println("TaskStore. TaskList. Error: ", err)
		return nil, err
	}

	return list, err
}

// TaskFindByID finds a Task by ID
func (st *Store) TaskFindByID(id string) *Task {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		From(st.taskTaskTableName).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).
		Select().
		Limit(1).
		ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var Task Task
	err := sqlscan.Get(context.Background(), st.db, &Task, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil
		}

		if sqlscan.NotFound(err) {
			return nil
		}

		log.Println("TaskSTore. TaskFindByID. Error: ", err)
		return nil
	}

	return &Task
}

// TaskUpdate creates a Task
func (st *Store) TaskUpdate(queue *Task) error {
	queue.UpdatedAt = time.Now()
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		Update(st.taskTaskTableName).
		Where(goqu.C("id").Eq(queue.ID)).
		Set(queue).
		ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		if st.debug {
			log.Println(err)
		}

		return err
	}

	return nil
}

// SqlCreateTaskTable returns a SQL string for creating the Task table
func (st *Store) SqlCreateTaskTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.taskTaskTableName + ` (
	  id             varchar(40) NOT NULL PRIMARY KEY,
	  status         varchar(40) NOT NULL,
	  definition_id  varchar(40) NOT NULL,
	  parameters     text        NOT NULL,
	  output         longtext    DEFAULT NULL,
	  details        longtext    DEFAULT NULL,
	  attempts       int         DEFAULT 0,
	  started_at     datetime    DEFAULT NULL,
	  completed_at   datetime    DEFAULT NULL,
	  created_at	 datetime,
	  updated_at	 datetime,	
	  deleted_at	 datetime    DEFAULT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.taskTaskTableName + `" (
	  "id"            varchar(40)    NOT NULL PRIMARY KEY,
	  "status"        varchar(40)    NOT NULL,
	  "definition_id" varchar(40)    NOT NULL,
	  "parameters"    varchar(40)    NOT NULL,
	  "output"        longtext       DEFAULT NULL,
	  "details"       longtext       DEFAULT NULL,
	  "attempts"      int            DEFAULT 0,
	  "started_at"    timestamptz(6) DEFAULT NULL,
	  "completed_at"  timestamptz(6) DEFAULT NULL,
	  "created_at"    timestamptz(6) NOT NULL,
	  "updated_at"    timestamptz(6) NOT NULL,
	  "deleted_at"    timestamptz(6) DEFAULT NULL
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.taskTaskTableName + `" (
	  "id"            varchar(40) NOT NULL PRIMARY KEY,
	  "status"        varchar(40) NOT NULL,
	  "definition_id" varchar(40) NOT NULL,
	  "parameters"    varchar(40) NOT NULL,
	  "output"        text  DEFAULT NULL,
	  "details"       text  DEFAULT NULL,
	  "attempts"      int   DEFAULT 0,
	  "started_at"    datetime  DEFAULT NULL,
	  "completed_at"  datetime  DEFAULT NULL,
	  "created_at"    datetime NOT NULL,
	  "updated_at"    datetime NOT NULL,
	  "deleted_at"    datetime DEFAULT NULL
	)
	`

	sql := "unsupported driver " + st.dbDriverName

	if st.dbDriverName == "mysql" {
		sql = sqlMysql
	}
	if st.dbDriverName == "postgres" {
		sql = sqlPostgres
	}
	if st.dbDriverName == "sqlite" {
		sql = sqlSqlite
	}

	return sql
}
