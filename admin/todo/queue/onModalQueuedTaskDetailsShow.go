package admin

// import (
// 	"project/config"

// 	"github.com/gouniverse/hb"
// )

// func (controller *queueManagerController) onModalQueuedTaskDetailsShow(queueID string) string {
// 	if queueID == "" {
// 		return hb.Div().Class("alert alert-danger").Text("queue id is required").ToHTML()
// 	}

// 	queue, err := config.TaskStore.QueueFindByID(queueID)

// 	if err != nil {
// 		config.Logger.Error("At taskadmin > onModalQueuedTaskDetailsShow", "error", err.Error())
// 		return hb.Div().Class("alert alert-danger").Text("Error retrieving queued task").ToHTML()
// 	}

// 	if queue == nil {
// 		return hb.Div().Class("alert alert-danger").Text("Queue task not found").ToHTML()
// 	}

// 	return controller.modalQueuedTaskDetails(queue.Details).ToHTML()
// }
