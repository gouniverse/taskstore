package taskstore

import "errors"

func QueueQuery() QueueQueryInterface {
	return &queueQuery{
		properties: make(map[string]interface{}),
	}
}

type queueQuery struct {
	properties map[string]interface{}
}

var _ QueueQueryInterface = (*queueQuery)(nil)

func (q *queueQuery) Validate() error {
	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("queue query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("queue query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("queue query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("queue query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("queue query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("queue query. offset cannot be negative")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("queue query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("queue query. status_in cannot be empty array")
	}

	if q.HasTaskID() && q.TaskID() == "" {
		return errors.New("queue query. task_id cannot be empty")
	}

	return nil
}

func (q *queueQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *queueQuery) SetColumns(columns []string) QueueQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *queueQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *queueQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *queueQuery) SetCountOnly(countOnly bool) QueueQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *queueQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *queueQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *queueQuery) SetCreatedAtGte(createdAtGte string) QueueQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *queueQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *queueQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *queueQuery) SetCreatedAtLte(createdAtLte string) QueueQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *queueQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *queueQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *queueQuery) SetID(id string) QueueQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *queueQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *queueQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *queueQuery) SetIDIn(idIn []string) QueueQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *queueQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *queueQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *queueQuery) SetLimit(limit int) QueueQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *queueQuery) HasTaskID() bool {
	return q.hasProperty("task_id")
}

func (q *queueQuery) TaskID() string {
	return q.properties["task_id"].(string)
}

func (q *queueQuery) SetTaskID(taskID string) QueueQueryInterface {
	q.properties["task_id"] = taskID
	return q
}

func (q *queueQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *queueQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *queueQuery) SetOffset(offset int) QueueQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *queueQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *queueQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *queueQuery) SetOrderBy(orderBy string) QueueQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *queueQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *queueQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *queueQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) QueueQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *queueQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *queueQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *queueQuery) SetSortOrder(sortOrder string) QueueQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *queueQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *queueQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *queueQuery) SetStatus(status string) QueueQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *queueQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *queueQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *queueQuery) SetStatusIn(statusIn []string) QueueQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *queueQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
