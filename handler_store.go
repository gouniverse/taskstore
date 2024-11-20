package taskstore

import "errors"

var errTaskMissing = errors.New("task not found")

func (store *Store) TaskHandlerAdd(taskHandler TaskHandlerInterface, createIfMissing bool) error {
	task, err := store.TaskFindByAlias(taskHandler.Alias())

	if err != nil {
		return err
	}

	if task == nil && !createIfMissing {
		return errTaskMissing
	}

	if task == nil && createIfMissing {
		alias := taskHandler.Alias()
		title := taskHandler.Title()
		description := taskHandler.Description()

		task := NewTask().
			SetStatus(TaskStatusActive).
			SetAlias(alias).
			SetTitle(title).
			SetDescription(description)

		err := store.TaskCreate(task)

		if err != nil {
			return err
		}
	}

	store.taskHandlers = append(store.taskHandlers, taskHandler)

	return nil
}

func (store *Store) TaskHandlerList() []TaskHandlerInterface {
	return store.taskHandlers
}
