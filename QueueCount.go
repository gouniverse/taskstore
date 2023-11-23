package taskstore

import (
	"context"
	"database/sql"
	"log"

	"github.com/doug-martin/goqu/v9"
	"github.com/georgysavva/scany/sqlscan"
)

func (st *Store) QueueCount(options QueueQueryOptions) (int64, error) {
	options.CountOnly = true

	q := st.queueQuery(options)

	sqlStr, _, errSql := q.Limit(1).
		Select(goqu.COUNT(goqu.Star()).As("count")).
		ToSQL()

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
