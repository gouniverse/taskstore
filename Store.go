package taskstore

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"reflect"
	"strings"
)

// Store defines a session store
type Store struct {
	taskDefinitionTableName string
	taskTaskTableName       string
	db                      *sql.DB
	dbDriverName            string
	automigrateEnabled      bool
	debug                   bool
}

// StoreOption options for the vault store
type StoreOption func(*Store)

// NewStore creates a new entity store
func NewStore(opts ...StoreOption) (*Store, error) {
	store := &Store{}
	for _, opt := range opts {
		opt(store)
	}

	if store.taskDefinitionTableName == "" {
		log.Panic("Task store: taskDefinitionTableName is required")
	}

	if store.taskTaskTableName == "" {
		log.Panic("Task store: taskTaskTableName is required")
	}

	if !store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// DriverName finds the driver name from database
func (st *Store) DriverName(db *sql.DB) string {
	dv := reflect.ValueOf(db.Driver())
	driverFullName := dv.Type().String()
	if strings.Contains(driverFullName, "mysql") {
		return "mysql"
	}
	if strings.Contains(driverFullName, "postgres") || strings.Contains(driverFullName, "pq") {
		return "postgres"
	}
	if strings.Contains(driverFullName, "sqlite") {
		return "sqlite"
	}
	if strings.Contains(driverFullName, "mssql") {
		return "mssql"
	}
	return driverFullName
}

// AutoMigrate migrates the tables
func (st *Store) AutoMigrate() error {
	sqlTask := st.SqlCreateDefinitionTable()

	if st.debug {
		log.Println(sqlTask)
	}

	_, errTask := st.db.Exec(sqlTask)
	if errTask != nil {
		log.Println(errTask)
		return errTask
	}

	sqlTaskTable := st.SqlCreateTaskTable()

	if st.debug {
		log.Println(sqlTaskTable)
	}

	_, errTaskTable := st.db.Exec(sqlTaskTable)
	if errTaskTable != nil {
		log.Println(errTaskTable)
		return errTaskTable
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debug bool) {
	st.debug = debug
}

// WithTaskTableDefinitionName sets the table name for the user
func WithDefinitionTableName(taskDefinitionTableName string) StoreOption {
	return func(s *Store) {
		s.taskDefinitionTableName = taskDefinitionTableName
	}
}

// WithLicenseTableName sets the table name for the email
func WithTaskTableName(taskTaskTableName string) StoreOption {
	return func(s *Store) {
		s.taskTaskTableName = taskTaskTableName
	}
}

// WithDb sets the database for the task store
func WithDb(db *sql.DB) StoreOption {
	return func(s *Store) {
		s.db = db
		s.dbDriverName = s.DriverName(s.db)
	}
}

// WithDebug prints the SQL queries
func WithDebug(debug bool) StoreOption {
	return func(s *Store) {
		s.debug = debug
	}
}

// EnqueueByAlias creates and enqueues a Task
func (st *Store) EnqueueByAlias(alias string, parameters map[string]interface{}) (*Task, error) {
	definition := st.DefinitionFindByAlias(alias)

	if definition == nil {
		return nil, errors.New("definition with alias '" + alias + "' not found")
	}

	parameters = prependTaskAlias(alias, parameters)

	parametersBytes, jsonErr := json.Marshal(parameters)

	if jsonErr != nil {
		return nil, errors.New("parameters json marshal error")
	}

	parametersStr := string(parametersBytes)

	queuedTask := Task{
		DefinitionID: definition.ID,
		Parameters:   parametersStr,
		Status:       TaskStatusQueued,
	}

	err := st.TaskCreate(&queuedTask)

	if err != nil {
		return &queuedTask, err
	}

	return &queuedTask, err
}

// prependTaskAlias prepends a task alias to the parameters so that its easy to distinguish
func prependTaskAlias(alias string, parameters map[string]interface{}) map[string]interface{} {
	copiedParameters := map[string]interface{}{
		"task_alias": alias,
	}
	for index, element := range parameters {
		copiedParameters[index] = element
	}

	return copiedParameters
}
