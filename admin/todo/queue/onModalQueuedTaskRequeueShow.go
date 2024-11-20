package admin

// import (
// 	"project/config"

// 	"github.com/gouniverse/hb"
// )

// func (controller *queueManagerController) onModalQueuedTaskRequeueShow(queueID string) string {
// 	queue, err := config.TaskStore.QueueFindByID(queueID)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At taskadmin > onModalQueuedTaskRequeueShow", err.Error())
// 		return hb.Div().Class("alert alert-danger").Text("Error retrieving queued task").ToHTML()
// 	}

// 	if queue == nil {
// 		return hb.Div().Class("alert alert-danger").Text("Queued task not found").ToHTML()
// 	}

// 	return controller.modalQueuedTaskRequeue(queue).ToHTML()
// }
