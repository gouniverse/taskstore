package admin

import (
	"net/http"
	"strings"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/utils"
)

func (controller *queueManagerController) onModalQueuedTaskDeleteSubmitted(r *http.Request) hb.TagInterface {
	queueID := strings.TrimSpace(utils.Req(r, "queue_id", ""))

	if queueID == "" {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Queued task ID is required"})
	}

	err := controller.store.QueueSoftDeleteByID(queueID)

	if err != nil {
		controller.logger.Error("At taskadmin > onModalQueuedTaskDeleteSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Queued task failed to be deleted"})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Queued task successfully deleted"})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 3000);`))
}
