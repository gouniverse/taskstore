package admin

// import (
// 	"net/http"
// 	"project/config"
// 	"strings"

// 	"github.com/gouniverse/hb"
// 	"github.com/gouniverse/utils"
// )

// func (controller *queueManagerController) onModalQueuedTaskDeleteSubmitted(r *http.Request) string {
// 	queueID := strings.TrimSpace(utils.Req(r, "queue_id", ""))

// 	if queueID == "" {
// 		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Queued task ID is required"}).ToHTML()
// 	}

// 	err := config.TaskStore.QueueSoftDeleteByID(queueID)

// 	if err != nil {
// 		config.Logger.Error("At taskadmin > onModalQueuedTaskDeleteSubmitted", "error", err.Error())
// 		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Queued task failed to be deleted"}).ToHTML()
// 	}

// 	return hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Queued task successfully deleted"}).ToHTML() +
// 		hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 3000);`).ToHTML()
// }
