package taskstore

import (
	"database/sql"
	"errors"
	"log"

	sb "github.com/gouniverse/sql"
)

// Store defines a session store
type Store struct {
	taskTableName      string
	queueTableName     string
	taskHandlers       []TaskHandlerInterface
	db                 *sql.DB
	dbDriverName       string
	automigrateEnabled bool
	debugEnabled       bool
}

// NewStoreOptions define the options for creating a new task store
type NewStoreOptions struct {
	TaskTableName      string
	QueueTableName     string
	DB                 *sql.DB
	DbDriverName       string
	AutomigrateEnabled bool
	DebugEnabled       bool
}

// NewStore creates a new task store
func NewStore(opts NewStoreOptions) (*Store, error) {
	store := &Store{
		taskTableName:      opts.TaskTableName,
		queueTableName:     opts.QueueTableName,
		automigrateEnabled: opts.AutomigrateEnabled,
		db:                 opts.DB,
		dbDriverName:       opts.DbDriverName,
		debugEnabled:       opts.DebugEnabled,
	}

	if store.taskTableName == "" {
		return nil, errors.New("task store: taskTableName is required")
	}

	if store.queueTableName == "" {
		return nil, errors.New("task store: queueTableName is required")
	}

	if store.db == nil {
		return nil, errors.New("task store: DB is required")
	}

	if store.dbDriverName == "" {
		store.dbDriverName = sb.DatabaseDriverName(store.db)
	}

	if store.automigrateEnabled {
		store.AutoMigrate()
	}

	return store, nil
}

// AutoMigrate migrates the tables
func (st *Store) AutoMigrate() error {
	sqlTaskTable := st.SqlCreateTaskTable()

	if st.debugEnabled {
		log.Println(sqlTaskTable)
	}

	_, errTask := st.db.Exec(sqlTaskTable)
	if errTask != nil {
		log.Println(errTask)
		return errTask
	}

	sqlQueueTable := st.SqlCreateQueueTable()

	if st.debugEnabled {
		log.Println(sqlQueueTable)
	}

	_, errQueue := st.db.Exec(sqlQueueTable)
	if errQueue != nil {
		log.Println(errQueue)
		return errQueue
	}

	return nil
}

// EnableDebug - enables the debug option
func (st *Store) EnableDebug(debugEnabled bool) *Store {
	st.debugEnabled = debugEnabled
	return st
}
