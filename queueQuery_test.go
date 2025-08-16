package taskstore

import (
	"testing"
)

func TestQueueQuery(t *testing.T) {
	query := QueueQuery()

	if query == nil {
		t.Fatal("QueueQuery: Expected query to be created, got nil")
	}

	// Test that it implements the interface
	var _ QueueQueryInterface = query
}

func TestQueueQuery_Validate(t *testing.T) {
	tests := []struct {
		name        string
		setupQuery  func() QueueQueryInterface
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid empty query",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery()
			},
			expectError: false,
		},
		{
			name: "valid query with all fields",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().
					SetCreatedAtGte("2023-01-01 00:00:00").
					SetCreatedAtLte("2023-12-31 23:59:59").
					SetID("test-id").
					SetIDIn([]string{"id1", "id2"}).
					SetLimit(10).
					SetOffset(0).
					SetStatus("queued").
					SetStatusIn([]string{"queued", "running"}).
					SetTaskID("task-123")
			},
			expectError: false,
		},
		{
			name: "empty created_at_gte",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetCreatedAtGte("")
			},
			expectError: true,
			errorMsg:    "queue query. created_at_gte cannot be empty",
		},
		{
			name: "empty created_at_lte",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetCreatedAtLte("")
			},
			expectError: true,
			errorMsg:    "queue query. created_at_lte cannot be empty",
		},
		{
			name: "empty id",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetID("")
			},
			expectError: true,
			errorMsg:    "queue query. id cannot be empty",
		},
		{
			name: "empty id_in array",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetIDIn([]string{})
			},
			expectError: true,
			errorMsg:    "queue query. id_in cannot be empty array",
		},
		{
			name: "negative limit",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetLimit(-1)
			},
			expectError: true,
			errorMsg:    "queue query. limit cannot be negative",
		},
		{
			name: "negative offset",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetOffset(-1)
			},
			expectError: true,
			errorMsg:    "queue query. offset cannot be negative",
		},
		{
			name: "empty status",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetStatus("")
			},
			expectError: true,
			errorMsg:    "queue query. status cannot be empty",
		},
		{
			name: "empty status_in array",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetStatusIn([]string{})
			},
			expectError: true,
			errorMsg:    "queue query. status_in cannot be empty array",
		},
		{
			name: "empty task_id",
			setupQuery: func() QueueQueryInterface {
				return QueueQuery().SetTaskID("")
			},
			expectError: true,
			errorMsg:    "queue query. task_id cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := tt.setupQuery()
			err := query.Validate()

			if tt.expectError {
				if err == nil {
					t.Errorf("Validate: Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Validate: Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Validate: Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestQueueQuery_Columns(t *testing.T) {
	query := QueueQuery()

	// Test default state
	columns := query.Columns()
	if len(columns) != 0 {
		t.Errorf("Columns: Expected empty slice, got %v", columns)
	}

	// Test setting columns
	testColumns := []string{"id", "task_id", "status"}
	result := query.SetColumns(testColumns)
	if result != query {
		t.Error("SetColumns: Expected method to return the same query instance")
	}
	
	retrievedColumns := query.Columns()
	if len(retrievedColumns) != len(testColumns) {
		t.Errorf("Columns: Expected %d columns, got %d", len(testColumns), len(retrievedColumns))
	}
	for i, col := range testColumns {
		if retrievedColumns[i] != col {
			t.Errorf("Columns: Expected column '%s' at index %d, got '%s'", col, i, retrievedColumns[i])
		}
	}
}

func TestQueueQuery_CountOnly(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasCountOnly() {
		t.Error("HasCountOnly: Expected false for new query")
	}
	if query.IsCountOnly() {
		t.Error("IsCountOnly: Expected false for new query")
	}

	// Test setting count only to true
	result := query.SetCountOnly(true)
	if result != query {
		t.Error("SetCountOnly: Expected method to return the same query instance")
	}
	if !query.HasCountOnly() {
		t.Error("HasCountOnly: Expected true after setting count only")
	}
	if !query.IsCountOnly() {
		t.Error("IsCountOnly: Expected true after setting count only to true")
	}

	// Test setting count only to false
	query.SetCountOnly(false)
	if !query.HasCountOnly() {
		t.Error("HasCountOnly: Expected true even when set to false")
	}
	if query.IsCountOnly() {
		t.Error("IsCountOnly: Expected false after setting count only to false")
	}
}

func TestQueueQuery_CreatedAtGte(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasCreatedAtGte() {
		t.Error("HasCreatedAtGte: Expected false for new query")
	}

	// Test setting created_at_gte
	testDate := "2023-01-01 00:00:00"
	result := query.SetCreatedAtGte(testDate)
	if result != query {
		t.Error("SetCreatedAtGte: Expected method to return the same query instance")
	}
	if !query.HasCreatedAtGte() {
		t.Error("HasCreatedAtGte: Expected true after setting created_at_gte")
	}
	if query.CreatedAtGte() != testDate {
		t.Errorf("CreatedAtGte: Expected '%s', got '%s'", testDate, query.CreatedAtGte())
	}
}

func TestQueueQuery_CreatedAtLte(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasCreatedAtLte() {
		t.Error("HasCreatedAtLte: Expected false for new query")
	}

	// Test setting created_at_lte
	testDate := "2023-12-31 23:59:59"
	result := query.SetCreatedAtLte(testDate)
	if result != query {
		t.Error("SetCreatedAtLte: Expected method to return the same query instance")
	}
	if !query.HasCreatedAtLte() {
		t.Error("HasCreatedAtLte: Expected true after setting created_at_lte")
	}
	if query.CreatedAtLte() != testDate {
		t.Errorf("CreatedAtLte: Expected '%s', got '%s'", testDate, query.CreatedAtLte())
	}
}

func TestQueueQuery_ID(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasID() {
		t.Error("HasID: Expected false for new query")
	}

	// Test setting ID
	testID := "test-queue-id-123"
	result := query.SetID(testID)
	if result != query {
		t.Error("SetID: Expected method to return the same query instance")
	}
	if !query.HasID() {
		t.Error("HasID: Expected true after setting ID")
	}
	if query.ID() != testID {
		t.Errorf("ID: Expected '%s', got '%s'", testID, query.ID())
	}
}

func TestQueueQuery_IDIn(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasIDIn() {
		t.Error("HasIDIn: Expected false for new query")
	}

	// Test setting ID in
	testIDs := []string{"queue1", "queue2", "queue3"}
	result := query.SetIDIn(testIDs)
	if result != query {
		t.Error("SetIDIn: Expected method to return the same query instance")
	}
	if !query.HasIDIn() {
		t.Error("HasIDIn: Expected true after setting ID in")
	}
	
	retrievedIDs := query.IDIn()
	if len(retrievedIDs) != len(testIDs) {
		t.Errorf("IDIn: Expected %d IDs, got %d", len(testIDs), len(retrievedIDs))
	}
	for i, id := range testIDs {
		if retrievedIDs[i] != id {
			t.Errorf("IDIn: Expected ID '%s' at index %d, got '%s'", id, i, retrievedIDs[i])
		}
	}
}

func TestQueueQuery_Limit(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasLimit() {
		t.Error("HasLimit: Expected false for new query")
	}

	// Test setting limit
	testLimit := 25
	result := query.SetLimit(testLimit)
	if result != query {
		t.Error("SetLimit: Expected method to return the same query instance")
	}
	if !query.HasLimit() {
		t.Error("HasLimit: Expected true after setting limit")
	}
	if query.Limit() != testLimit {
		t.Errorf("Limit: Expected %d, got %d", testLimit, query.Limit())
	}
}

func TestQueueQuery_TaskID(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasTaskID() {
		t.Error("HasTaskID: Expected false for new query")
	}

	// Test setting task ID
	testTaskID := "task-456"
	result := query.SetTaskID(testTaskID)
	if result != query {
		t.Error("SetTaskID: Expected method to return the same query instance")
	}
	if !query.HasTaskID() {
		t.Error("HasTaskID: Expected true after setting task ID")
	}
	if query.TaskID() != testTaskID {
		t.Errorf("TaskID: Expected '%s', got '%s'", testTaskID, query.TaskID())
	}
}

func TestQueueQuery_Offset(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasOffset() {
		t.Error("HasOffset: Expected false for new query")
	}

	// Test setting offset
	testOffset := 15
	result := query.SetOffset(testOffset)
	if result != query {
		t.Error("SetOffset: Expected method to return the same query instance")
	}
	if !query.HasOffset() {
		t.Error("HasOffset: Expected true after setting offset")
	}
	if query.Offset() != testOffset {
		t.Errorf("Offset: Expected %d, got %d", testOffset, query.Offset())
	}
}

func TestQueueQuery_OrderBy(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasOrderBy() {
		t.Error("HasOrderBy: Expected false for new query")
	}

	// Test setting order by
	testOrderBy := "started_at"
	result := query.SetOrderBy(testOrderBy)
	if result != query {
		t.Error("SetOrderBy: Expected method to return the same query instance")
	}
	if !query.HasOrderBy() {
		t.Error("HasOrderBy: Expected true after setting order by")
	}
	if query.OrderBy() != testOrderBy {
		t.Errorf("OrderBy: Expected '%s', got '%s'", testOrderBy, query.OrderBy())
	}
}

func TestQueueQuery_SoftDeletedIncluded(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasSoftDeletedIncluded() {
		t.Error("HasSoftDeletedIncluded: Expected false for new query")
	}
	if query.SoftDeletedIncluded() {
		t.Error("SoftDeletedIncluded: Expected false for new query")
	}

	// Test setting soft deleted included to true
	result := query.SetSoftDeletedIncluded(true)
	if result != query {
		t.Error("SetSoftDeletedIncluded: Expected method to return the same query instance")
	}
	if !query.HasSoftDeletedIncluded() {
		t.Error("HasSoftDeletedIncluded: Expected true after setting soft deleted included")
	}
	if !query.SoftDeletedIncluded() {
		t.Error("SoftDeletedIncluded: Expected true after setting to true")
	}

	// Test setting soft deleted included to false
	query.SetSoftDeletedIncluded(false)
	if !query.HasSoftDeletedIncluded() {
		t.Error("HasSoftDeletedIncluded: Expected true even when set to false")
	}
	if query.SoftDeletedIncluded() {
		t.Error("SoftDeletedIncluded: Expected false after setting to false")
	}
}

func TestQueueQuery_SortOrder(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasSortOrder() {
		t.Error("HasSortOrder: Expected false for new query")
	}

	// Test setting sort order
	testSortOrder := "ASC"
	result := query.SetSortOrder(testSortOrder)
	if result != query {
		t.Error("SetSortOrder: Expected method to return the same query instance")
	}
	if !query.HasSortOrder() {
		t.Error("HasSortOrder: Expected true after setting sort order")
	}
	if query.SortOrder() != testSortOrder {
		t.Errorf("SortOrder: Expected '%s', got '%s'", testSortOrder, query.SortOrder())
	}
}

func TestQueueQuery_Status(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasStatus() {
		t.Error("HasStatus: Expected false for new query")
	}

	// Test setting status
	testStatus := "running"
	result := query.SetStatus(testStatus)
	if result != query {
		t.Error("SetStatus: Expected method to return the same query instance")
	}
	if !query.HasStatus() {
		t.Error("HasStatus: Expected true after setting status")
	}
	if query.Status() != testStatus {
		t.Errorf("Status: Expected '%s', got '%s'", testStatus, query.Status())
	}
}

func TestQueueQuery_StatusIn(t *testing.T) {
	query := QueueQuery()

	// Test default state
	if query.HasStatusIn() {
		t.Error("HasStatusIn: Expected false for new query")
	}

	// Test setting status in
	testStatuses := []string{"queued", "running", "success"}
	result := query.SetStatusIn(testStatuses)
	if result != query {
		t.Error("SetStatusIn: Expected method to return the same query instance")
	}
	if !query.HasStatusIn() {
		t.Error("HasStatusIn: Expected true after setting status in")
	}
	
	retrievedStatuses := query.StatusIn()
	if len(retrievedStatuses) != len(testStatuses) {
		t.Errorf("StatusIn: Expected %d statuses, got %d", len(testStatuses), len(retrievedStatuses))
	}
	for i, status := range testStatuses {
		if retrievedStatuses[i] != status {
			t.Errorf("StatusIn: Expected status '%s' at index %d, got '%s'", status, i, retrievedStatuses[i])
		}
	}
}

func TestQueueQuery_ChainedSetters(t *testing.T) {
	query := QueueQuery()

	// Test that all setters can be chained
	result := query.
		SetColumns([]string{"id", "task_id"}).
		SetCountOnly(true).
		SetCreatedAtGte("2023-01-01 00:00:00").
		SetCreatedAtLte("2023-12-31 23:59:59").
		SetID("test-id").
		SetIDIn([]string{"id1", "id2"}).
		SetLimit(20).
		SetTaskID("task-789").
		SetOffset(10).
		SetOrderBy("created_at").
		SetSoftDeletedIncluded(true).
		SetSortOrder("DESC").
		SetStatus("queued").
		SetStatusIn([]string{"queued", "running"})

	if result != query {
		t.Error("ChainedSetters: Expected all setters to return the same query instance for chaining")
	}

	// Verify all values were set correctly
	if len(query.Columns()) != 2 {
		t.Error("ChainedSetters: Columns not set correctly")
	}
	if !query.IsCountOnly() {
		t.Error("ChainedSetters: CountOnly not set correctly")
	}
	if query.CreatedAtGte() != "2023-01-01 00:00:00" {
		t.Error("ChainedSetters: CreatedAtGte not set correctly")
	}
	if query.CreatedAtLte() != "2023-12-31 23:59:59" {
		t.Error("ChainedSetters: CreatedAtLte not set correctly")
	}
	if query.ID() != "test-id" {
		t.Error("ChainedSetters: ID not set correctly")
	}
	if len(query.IDIn()) != 2 {
		t.Error("ChainedSetters: IDIn not set correctly")
	}
	if query.Limit() != 20 {
		t.Error("ChainedSetters: Limit not set correctly")
	}
	if query.TaskID() != "task-789" {
		t.Error("ChainedSetters: TaskID not set correctly")
	}
	if query.Offset() != 10 {
		t.Error("ChainedSetters: Offset not set correctly")
	}
	if query.OrderBy() != "created_at" {
		t.Error("ChainedSetters: OrderBy not set correctly")
	}
	if !query.SoftDeletedIncluded() {
		t.Error("ChainedSetters: SoftDeletedIncluded not set correctly")
	}
	if query.SortOrder() != "DESC" {
		t.Error("ChainedSetters: SortOrder not set correctly")
	}
	if query.Status() != "queued" {
		t.Error("ChainedSetters: Status not set correctly")
	}
	if len(query.StatusIn()) != 2 {
		t.Error("ChainedSetters: StatusIn not set correctly")
	}
}