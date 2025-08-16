package taskstore

import "errors"

func TaskQuery() TaskQueryInterface {
	return &taskQuery{
		properties: make(map[string]interface{}),
	}
}

type taskQuery struct {
	properties map[string]interface{}
}

var _ TaskQueryInterface = (*taskQuery)(nil)

func (q *taskQuery) Validate() error {
	if q.HasAlias() && q.Alias() == "" {
		return errors.New("task query. alias cannot be empty")
	}

	if q.HasCreatedAtGte() && q.CreatedAtGte() == "" {
		return errors.New("task query. created_at_gte cannot be empty")
	}

	if q.HasCreatedAtLte() && q.CreatedAtLte() == "" {
		return errors.New("task query. created_at_lte cannot be empty")
	}

	if q.HasID() && q.ID() == "" {
		return errors.New("task query. id cannot be empty")
	}

	if q.HasIDIn() && len(q.IDIn()) < 1 {
		return errors.New("task query. id_in cannot be empty array")
	}

	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("task query. limit cannot be negative")
	}

	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("task query. offset cannot be negative")
	}

	if q.HasStatus() && q.Status() == "" {
		return errors.New("task query. status cannot be empty")
	}

	if q.HasStatusIn() && len(q.StatusIn()) < 1 {
		return errors.New("task query. status_in cannot be empty array")
	}
	
	return nil
}

func (q *taskQuery) HasAlias() bool {
	return q.hasProperty("alias")
}

func (q *taskQuery) Alias() string {
	if !q.hasProperty("alias") {
		return ""
	}

	return q.properties["alias"].(string)
}

func (q *taskQuery) SetAlias(alias string) TaskQueryInterface {
	q.properties["alias"] = alias
	return q
}

func (q *taskQuery) Columns() []string {
	if !q.hasProperty("columns") {
		return []string{}
	}

	return q.properties["columns"].([]string)
}

func (q *taskQuery) SetColumns(columns []string) TaskQueryInterface {
	q.properties["columns"] = columns
	return q
}

func (q *taskQuery) HasCountOnly() bool {
	return q.hasProperty("count_only")
}

func (q *taskQuery) IsCountOnly() bool {
	if q.HasCountOnly() {
		return q.properties["count_only"].(bool)
	}

	return false
}

func (q *taskQuery) SetCountOnly(countOnly bool) TaskQueryInterface {
	q.properties["count_only"] = countOnly
	return q
}

func (q *taskQuery) HasCreatedAtGte() bool {
	return q.hasProperty("created_at_gte")
}

func (q *taskQuery) CreatedAtGte() string {
	return q.properties["created_at_gte"].(string)
}

func (q *taskQuery) SetCreatedAtGte(createdAtGte string) TaskQueryInterface {
	q.properties["created_at_gte"] = createdAtGte
	return q
}

func (q *taskQuery) HasCreatedAtLte() bool {
	return q.hasProperty("created_at_lte")
}

func (q *taskQuery) CreatedAtLte() string {
	return q.properties["created_at_lte"].(string)
}

func (q *taskQuery) SetCreatedAtLte(createdAtLte string) TaskQueryInterface {
	q.properties["created_at_lte"] = createdAtLte
	return q
}

func (q *taskQuery) HasID() bool {
	return q.hasProperty("id")
}

func (q *taskQuery) ID() string {
	return q.properties["id"].(string)
}

func (q *taskQuery) SetID(id string) TaskQueryInterface {
	q.properties["id"] = id
	return q
}

func (q *taskQuery) HasIDIn() bool {
	return q.hasProperty("id_in")
}

func (q *taskQuery) IDIn() []string {
	return q.properties["id_in"].([]string)
}

func (q *taskQuery) SetIDIn(idIn []string) TaskQueryInterface {
	q.properties["id_in"] = idIn
	return q
}

func (q *taskQuery) HasLimit() bool {
	return q.hasProperty("limit")
}

func (q *taskQuery) Limit() int {
	return q.properties["limit"].(int)
}

func (q *taskQuery) SetLimit(limit int) TaskQueryInterface {
	q.properties["limit"] = limit
	return q
}

func (q *taskQuery) HasOffset() bool {
	return q.hasProperty("offset")
}

func (q *taskQuery) Offset() int {
	return q.properties["offset"].(int)
}

func (q *taskQuery) SetOffset(offset int) TaskQueryInterface {
	q.properties["offset"] = offset
	return q
}

func (q *taskQuery) HasOrderBy() bool {
	return q.hasProperty("order_by")
}

func (q *taskQuery) OrderBy() string {
	return q.properties["order_by"].(string)
}

func (q *taskQuery) SetOrderBy(orderBy string) TaskQueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

func (q *taskQuery) HasSoftDeletedIncluded() bool {
	return q.hasProperty("soft_delete_included")
}

func (q *taskQuery) SoftDeletedIncluded() bool {
	if !q.HasSoftDeletedIncluded() {
		return false
	}
	return q.properties["soft_delete_included"].(bool)
}

func (q *taskQuery) SetSoftDeletedIncluded(softDeleteIncluded bool) TaskQueryInterface {
	q.properties["soft_delete_included"] = softDeleteIncluded
	return q
}

func (q *taskQuery) HasSortOrder() bool {
	return q.hasProperty("sort_order")
}

func (q *taskQuery) SortOrder() string {
	return q.properties["sort_order"].(string)
}

func (q *taskQuery) SetSortOrder(sortOrder string) TaskQueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}

func (q *taskQuery) HasStatus() bool {
	return q.hasProperty("status")
}

func (q *taskQuery) Status() string {
	return q.properties["status"].(string)
}

func (q *taskQuery) SetStatus(status string) TaskQueryInterface {
	q.properties["status"] = status
	return q
}

func (q *taskQuery) HasStatusIn() bool {
	return q.hasProperty("status_in")
}

func (q *taskQuery) StatusIn() []string {
	return q.properties["status_in"].([]string)
}

func (q *taskQuery) SetStatusIn(statusIn []string) TaskQueryInterface {
	q.properties["status_in"] = statusIn
	return q
}

func (q *taskQuery) hasProperty(key string) bool {
	return q.properties[key] != nil
}
