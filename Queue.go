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

// QueueFail fails a queued task
func (st *Store) QueueFail(queue *Queue) error {
	completedAt := time.Now()
	queue.CompletedAt = &completedAt
	queue.Status = QueueStatusFailed
	return st.QueueUpdate(queue)
}

// QueueSuccess completes a queued task  successfully
func (st *Store) QueueSuccess(queue *Queue) error {
	completedAt := time.Now()
	queue.CompletedAt = &completedAt
	queue.Status = QueueStatusSuccess
	return st.QueueUpdate(queue)
}

// QueueCreate creates a queued task
func (st *Store) QueueCreate(queue *Queue) error {
	if queue.ID == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		queue.ID = uid.MicroUid()
	}
	queue.CreatedAt = time.Now()
	queue.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.queueTableName).Rows(queue).ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return err
	}

	return nil
}

type QueueQueryOptions struct {
	Offset    int64
	Status    string
	Limit     int
	SortBy    string
	SortOrder string
	CountOnly bool
}

func (st *Store) queueQuery(options QueueQueryOptions) *goqu.SelectDataset {
	q := goqu.Dialect(st.dbDriverName).From(st.queueTableName)

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

	if options.Offset != 0 && !options.CountOnly {
		q = q.Offset(uint(options.Offset))
	}

	return q
}

func (st *Store) QueueCount(options QueueQueryOptions) (int64, error) {
	options.CountOnly = true

	q := st.queueQuery(options)

	sqlStr, _, errSql := q.Limit(1).Select(goqu.COUNT(goqu.Star()).As("count")).ToSQL()

	if errSql != nil {
		return -1, nil
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	type Count struct {
		Count int64 `db:"count"`
	}
	count := []Count{}
	err := sqlscan.Select(context.Background(), st.db, &count, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return 0, nil
		}

		if sqlscan.NotFound(err) {
			return 0, nil
		}

		log.Println("QueueStore. QueueList. Error: ", err)
		return 0, err
	}

	if len(count) == 0 {
		return 0, nil
	}

	return count[0].Count, err
}

func (st *Store) QueueList(options QueueQueryOptions) ([]Queue, error) {
	q := st.queueQuery(options)

	sqlStr, _, _ := q.Select().ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	list := []Queue{}
	err := sqlscan.Select(context.Background(), st.db, &list, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		log.Println("QueueStore. QueueList. Error: ", err)
		return nil, err
	}

	return list, err
}

// QueueFindByID finds a Queue by ID
func (st *Store) QueueFindByID(id string) (*Queue, error) {
	sqlStr, _, err := goqu.Dialect(st.dbDriverName).
		From(st.queueTableName).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).
		Select().
		Limit(1).
		ToSQL()

	if err != nil {
		return nil, err
	}

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	var Queue Queue
	err = sqlscan.Get(context.Background(), st.db, &Queue, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		log.Println("QueueStore. QueueFindByID. Error: ", err)
		return nil, err
	}

	return &Queue, nil
}

// QueueUpdate creates a Queue
func (st *Store) QueueUpdate(queue *Queue) error {
	queue.UpdatedAt = time.Now()
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		Update(st.queueTableName).
		Where(goqu.C("id").Eq(queue.ID)).
		Set(queue).
		ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		if st.debugEnabled {
			log.Println(err)
		}

		return err
	}

	return nil
}
