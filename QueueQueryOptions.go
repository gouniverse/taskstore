package taskstore

type QueueQueryOptions struct {
	ID           string
	TaskID       string
	Status       string
	CreatedAtGte string
	UpdatedAtLte string
	Offset       int64
	Limit        int
	SortBy       string
	SortOrder    string
	CountOnly    bool
}
