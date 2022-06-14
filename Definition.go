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
	DefinitionStatusActive   = "active"
	DefinitionStatusCanceled = "canceled"
)

// Definition type
type Definition struct {
	ID          string     `json:"id" db:"id"`                   // varchar(40)  primary_key
	Status      string     `json:"status" db:"status"`           // varchar(40)  NOT NULL
	Alias       string     `json:"alias" db:"alias"`             // varchar(40)  NOT NULL
	Title       string     `json:"title" db:"title"`             // varchar(255) NOT NULL
	Description *string    `json:"description" db:"description"` // text         DEFAULT NULL
	Memo        *string    `json:"memo" db:"memo"`               // text         DEFAULT NULL
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`   // datetime     NOT NULL
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`   // datetime     NOT NULL
	DeletedAt   *time.Time `json:"deleted_at" db:"deleted_at"`   // datetime     DEFAULT NULL
}

// DefinitionCreate creates a Definition
func (st *Store) DefinitionCreate(Definition *Definition) (bool, error) {
	if Definition.ID == "" {
		Definition.ID = uid.MicroUid()
	}
	Definition.CreatedAt = time.Now()
	Definition.UpdatedAt = time.Now()

	var sqlStr string
	sqlStr, _, _ = goqu.Dialect(st.dbDriverName).Insert(st.taskDefinitionTableName).Rows(Definition).ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	_, err := st.db.Exec(sqlStr)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (st *Store) DefinitionList(options map[string]string) ([]Definition, error) {
	status, statusExists := options["status"]
	if !statusExists {
		status = ""
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taskDefinitionTableName)

	if status != "" {
		q = q.Where(goqu.C("status").Eq(status))
	}

	q = q.Where(goqu.C("deleted_at").IsNull())
	sqlStr, _, _ := q.Select().ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	list := []Definition{}
	err := sqlscan.Select(context.Background(), st.db, &list, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, nil
		}
		// log.Fatal("Failed to execute query: ", err)
		return nil, err
	}

	return list, err
}

// DefinitionFindByAlias finds a Definition by Alias
func (st *Store) DefinitionFindByAlias(alias string) *Definition {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.taskDefinitionTableName).Where(goqu.C("alias").Eq(alias), goqu.C("deleted_at").IsNull()).Select().Limit(1).ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var Definition Definition
	err := sqlscan.Get(context.Background(), st.db, &Definition, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	return &Definition
}

// DefinitionFindByID finds a Definition by ID
func (st *Store) DefinitionFindByID(id string) *Definition {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).From(st.taskDefinitionTableName).Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).Select().Limit(1).ToSQL()

	if st.debug {
		log.Println(sqlStr)
	}

	var Definition Definition
	err := sqlscan.Get(context.Background(), st.db, &Definition, sqlStr)

	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil
		}
		log.Fatal("Failed to execute query: ", err)
		return nil
	}

	return &Definition
}

// DefinitionUpdate creates a Definition
func (st *Store) DefinitionUpdate(Definition *Definition) bool {

	// result := st.db.Table(st.DefinitionTableName).Save(&Definition)

	// if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	// 	return false
	// }

	return true
}

// SqlCreateDefinitionTable returns a SQL string for creating the Definition table
func (st *Store) SqlCreateDefinitionTable() string {
	sqlMysql := `
	CREATE TABLE IF NOT EXISTS ` + st.taskDefinitionTableName + ` (
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
	CREATE TABLE IF NOT EXISTS "` + st.taskDefinitionTableName + `" (
	  "id"          varchar(40)    NOT NULL PRIMARY KEY,
	  "status"      varchar(40)    NOT NULL,
	  "alias"       varchar(40)    NOT NULL,
	  "title"       varchar(255)   NOT NULL,
	  "description" text,          DEFAULT NULL
	  "memo"        text,          DEFAULT NULL
	  "created_at"  timestamptz(6) NOT NULL,
	  "updated_at"  timestamptz(6) NOT NULL,
	  "deleted_at"  timestamptz(6) DEFAULT NULL
	)
	`

	sqlSqlite := `
	CREATE TABLE IF NOT EXISTS "` + st.taskDefinitionTableName + `" (
	  "id"          varchar(40)  NOT NULL PRIMARY KEY,
	  "status"      varchar(40)  NOT NULL,
	  "alias"       varchar(40)  NOT NULL,
	  "title"       varchar(255) NOT NULL,
	  "description" text,
	  "memo"        text,
	  "created_at"  datetime     NOT NULL,
	  "updated_at"  datetime     NOT NULL,
	  "deleted_at"  datetime     DEFAULT NULL
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
