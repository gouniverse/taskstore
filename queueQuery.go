package taskstore

import "github.com/doug-martin/goqu/v9"

type QueueQueryOptions struct {
	TaskID    string
	Status    string
	Offset    int64
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

	if options.TaskID != "" {
		q = q.Where(goqu.C("task_id").Eq(options.TaskID))
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
