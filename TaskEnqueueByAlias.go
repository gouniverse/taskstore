package taskstore

import (
	"encoding/json"
	"errors"
)

// TaskEnqueueByAlias finds a task by its alias and appends it to the queue
func (st *Store) TaskEnqueueByAlias(taskAlias string, parameters map[string]interface{}) (*Queue, error) {
	task := st.TaskFindByAlias(taskAlias)

	if task == nil {
		return nil, errors.New("task with alias '" + taskAlias + "' not found")
	}

	parameters = queuePrependTaskAliasToParameters(taskAlias, parameters)

	parametersBytes, jsonErr := json.Marshal(parameters)

	if jsonErr != nil {
		return nil, errors.New("parameters json marshal error")
	}

	parametersStr := string(parametersBytes)

	queuedTask := Queue{
		TaskID:     task.ID,
		Parameters: parametersStr,
		Status:     QueueStatusQueued,
	}

	err := st.QueueCreate(&queuedTask)

	if err != nil {
		return &queuedTask, err
	}

	return &queuedTask, err
}

// queuePrependTaskAliasToParameters prepends a task alias to the queue parameters so that its easy to distinguish
func queuePrependTaskAliasToParameters(alias string, parameters map[string]interface{}) map[string]interface{} {
	copiedParameters := map[string]interface{}{
		"task_alias": alias,
	}
	for index, element := range parameters {
		copiedParameters[index] = element
	}

	return copiedParameters
}
