package taskstore

import "errors"

var errTaskMissing = errors.New("task not found")

func (store *Store) TaskHandlerAdd(handler TaskHandlerInterface, createIfMissing bool) error {
	task := store.TaskFindByAlias(handler.Alias())

	if task == nil && !createIfMissing {
		return errTaskMissing
	}

	if task == nil && createIfMissing {
		alias := handler.Alias()
		title := handler.Title()
		description := handler.Description()

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

	return nil
}
