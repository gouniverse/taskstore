package taskstore

type QueueQueryOptions struct {
	TaskID            string
	Status            string
	CreatedAtLessThan string
	UpdatedAtLessThan string
	Offset            int64
	Limit             int
	SortBy            string
	SortOrder         string
	CountOnly         bool
}
