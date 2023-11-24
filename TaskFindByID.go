package taskstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
)

// TaskFindByID finds a Task by ID
func (st *Store) TaskFindByID(id string) *Task {
	sqlStr, _, _ := goqu.Dialect(st.dbDriverName).
		From(st.taskTableName).
		Where(goqu.C("id").Eq(id), goqu.C("deleted_at").IsNull()).
		Select().
		Limit(1).
		ToSQL()

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
