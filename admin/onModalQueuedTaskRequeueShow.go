package admin

import (
	"net/http"

	"github.com/gouniverse/hb"
)

func (controller *queueManagerController) onModalQueuedTaskRequeueShow(r *http.Request, queueID string) hb.TagInterface {
	queue, err := controller.store.QueueFindByID(queueID)

	if err != nil {
		controller.logger.Error("At taskadmin > onModalQueuedTaskRequeueShow", "error", err.Error())
		return hb.Div().Class("alert alert-danger").Text("Error retrieving queued task")
	}

	if queue == nil {
		return hb.Div().Class("alert alert-danger").Text("Queued task not found")
	}

	return controller.modalQueuedTaskRequeue(r, queue)
}
