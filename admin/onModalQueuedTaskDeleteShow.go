package admin

import (
	"net/http"
	"strings"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func (controller *queueManagerController) onModalQueuedTaskDeleteShow(r *http.Request) hb.TagInterface {
	queueID := strings.TrimSpace(utils.Req(r, "queue_id", ""))

	if queueID == "" {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Queued task ID is required"})
	}

	queue, err := controller.store.QueueFindByID(queueID)

	if err != nil {
		controller.logger.Error("At taskadmin > onModalQueuedTaskDeleteShow", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Error retrieving queued task"})
	}

	if queue == nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Queued task not found"})
	}

	return controller.modalQueuedTaskDelete(r, queueID)
}
