package taskstore

type TaskQueryOptions struct {
	Alias                string
	ID                   string
	IDIn                 []string
	Status               string
	StatusIn             []string
	CreatedAtGreaterThan string
	UpdatedAtGreaterThan string
	CreatedAtLessThan    string
	UpdatedAtLessThan    string
	Offset               int
	Limit                int
	SortOrder            string
	OrderBy              string
	CountOnly            bool
	WithDeleted          bool
}
