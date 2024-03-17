package taskstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
)

func (st *Store) TaskList(options map[string]string) ([]Task, error) {
	status, statusExists := options[COLUMN_STATUS]
	if !statusExists {
		status = ""
	}

	q := goqu.Dialect(st.dbDriverName).From(st.taskTableName)

	if status != "" {
		q = q.Where(goqu.C(COLUMN_STATUS).Eq(status))
	}

	q = q.Where(goqu.C(COLUMN_DELETED_AT).IsNull())
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
