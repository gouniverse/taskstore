package taskstore

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/gouniverse/uid"
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

// TaskCreate creates a Task
func (st *Store) TaskCreate(Task *Task) (bool, error) {
	if Task.ID == "" {
		time.Sleep(1 * time.Millisecond) // !!! important
		Task.ID = uid.MicroUid()
	}
	Task.CreatedAt = time.Now()
	Task.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.taskTableName).Rows(Task).ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (st *Store) TaskList(options map[string]string) ([]Task, error) {
	status, statusExists := options["status"]
	if !statusExists {
		status = ""
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taskTableName)

	if status != "" {
		q = q.Where(goqu.C("status").Eq(status))
	}

	q = q.Where(goqu.C("deleted_at").IsNull())
	sqlStr, _, _ := q.Select().ToSQL()

	if st.debugEnabled {
		log.Println(sqlStr)
	}

	list := []Task{}
	err := sqlscan.Select(context.Background(), st.db, &list, sqlStr)

	if err != nil {
		if err == sql.ErrNoRows {
			// sqlscan does not use this anymore
			return nil, nil
		}

		if sqlscan.NotFound(err) {
			return nil, nil
		}

		log.Println("TaskSTore. Error: ", err)
		return nil, err
	}

	return list, err
}

// TaskFindByAlias finds a Task by Alias
func (st *Store) TaskFindByAlias(alias string) *Task {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.taskTableName).Where(goqu.C("alias").Eq(alias), goqu.C("deleted_at").IsNull()).Select().Limit(1).ToSQL()

	if st.debugEnabled {
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

		log.Println("TaskStore. TaskFindByAlias. Error: ", err)
		return nil
	}

	return &Task
}

// TaskFindByID finds a Task by ID
func (st *Store) TaskFindByID(id string) *Task {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.taskTableName).Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).Select().Limit(1).ToSQL()

	if st.debugEnabled {
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

		log.Fatal("TaskStore. TaskFindByID. Error: ", err)
		return nil
	}

	return &Task
}

// TaskUpdate creates a Task
func (st *Store) TaskUpdate(Task *Task) bool {

	// result := st.db.Table(st.TaskTableName).Save(&Task)

	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	return false
	// }

	return true
}

// SqlCreateTaskTable returns a SQL string for creating the Task table
func (st *Store) SqlCreateTaskTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.taskTableName + ` (
	  id             varchar(40)  NOT NULL PRIMARY KEY,
	  status         varchar(40)  NOT NULL,
	  alias          varchar(100) NOT NULL,
	  title          varchar(255) NOT NULL,
	  description    text         DEFAULT NULL,
	  memo           text         DEFAULT NULL,
	  created_at	 datetime     NOT NULL,
	  updated_at	 datetime 	  NOT NULL,
	  deleted_at	 datetime     DEFAULT NULL
	);
	`

	sqlPostgres := `
	CREATE TABLE IF NOT EXISTS "` + st.taskTableName + `" (
	  "id"          varchar(40)    NOT NULL PRIMARY KEY,
	  "status"      varchar(40)    NOT NULL,
	  "alias"       varchar(40)    NOT NULL,
	  "title"       varchar(255)   NOT NULL,
	  "description" text,          DEFAULT NULL
	  "memo"        text,          DEFAULT NULL
	  "created_at"  timestamptz(6) NOT NULL,
	  "updated_at"  timestamptz(6) NOT NULL,
	  "deleted_at"  timestamptz(6) DEFAULT NULL
	);
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.taskTableName + `" (
	  "id"          varchar(40)  NOT NULL PRIMARY KEY,
	  "status"      varchar(40)  NOT NULL,
	  "alias"       varchar(40)  NOT NULL,
	  "title"       varchar(255) NOT NULL,
	  "description" text,
	  "memo"        text,
	  "created_at"  datetime     NOT NULL,
	  "updated_at"  datetime     NOT NULL,
	  "deleted_at"  datetime     DEFAULT NULL
	);
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
