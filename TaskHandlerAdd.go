package taskstore

import "errors"

var errTaskMissing = errors.New("task not found")

func (store *Store) TaskHandlerAdd(taskHandler TaskHandlerInterface, createIfMissing bool) error {
	task := store.TaskFindByAlias(taskHandler.Alias())

	if task == nil && !createIfMissing {
		return errTaskMissing
	}

	if task == nil && createIfMissing {
		alias := taskHandler.Alias()
		title := taskHandler.Title()
		description := taskHandler.Description()

		task := Task{
			Status:      TaskStatusActive,
			Alias:       alias,
			Title:       title,
			Description: description,
		}
		_, err := store.TaskCreate(&task)
		if err != nil {
			return err
		}
	}

	store.taskHandlers = append(store.taskHandlers, taskHandler)

	return nil
}
