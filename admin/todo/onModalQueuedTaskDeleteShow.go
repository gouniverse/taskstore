package admin

// import (
// 	"net/http"
// 	"project/config"
// 	"strings"

// 	"github.com/gouniverse/hb"
// 	"github.com/gouniverse/utils"
// )

// func (controller *queueManagerController) onModalQueuedTaskDeleteShow(r *http.Request) string {
// 	queueID := strings.TrimSpace(utils.Req(r, "queue_id", ""))

// 	if queueID == "" {
// 		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Queued task ID is required"}).ToHTML()
// 	}

// 	queue, err := config.TaskStore.QueueFindByID(queueID)

// 	if err != nil {
// 		config.Logger.Error("At taskadmin > onModalQueuedTaskDeleteShow", "error", err.Error())
// 		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Error retrieving queued task"}).ToHTML()
// 	}

// 	if queue == nil {
// 		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Queued task not found"}).ToHTML()
// 	}

// 	return controller.modalQueuedTaskDelete(r, queueID).ToHTML()
// }
