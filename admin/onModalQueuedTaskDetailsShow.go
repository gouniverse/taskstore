package admin

import (
	"github.com/gouniverse/hb"
)

func (c *queueManagerController) onModalQueuedTaskDetailsShow(queueID string) *hb.Tag {
	if queueID == "" {
		return hb.Div().Class("alert alert-danger").Text("queue id is required")
	}

	queue, err := c.store.QueueFindByID(queueID)

	if err != nil {
		c.logger.Error("At taskadmin > onModalQueuedTaskDetailsShow", "error", err.Error())
		return hb.Div().Class("alert alert-danger").Text("Error retrieving queued task")
	}

	if queue == nil {
		return hb.Div().Class("alert alert-danger").Text("Queue task not found")
	}

	return c.modalQueuedTaskDetails(queue.Details())
}
