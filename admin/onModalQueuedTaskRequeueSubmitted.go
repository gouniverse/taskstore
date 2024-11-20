package admin

import (
	"net/http"
	"strings"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/utils"
)

func (controller *queueManagerController) onModalQueuedTaskRequeueSubmitted(r *http.Request) hb.TagInterface {
	taskID := strings.TrimSpace(utils.Req(r, "task_id", ""))
	taskParameters := strings.TrimSpace(utils.Req(r, "task_parameters", ""))

	if taskID == "" {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Task is required"})
	}

	if taskParameters == "" {
		taskParameters = "{}"
	}

	if !utils.IsJSON(taskParameters) {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Task Parameters is not valid JSON"})
	}

	task, err := controller.store.TaskFindByID(taskID)

	if err != nil {
		controller.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	if task == nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Task not found"})
	}

	taskParametersAny, err := utils.FromJSON(taskParameters, map[string]interface{}{})

	if err != nil {
		controller.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.Div().Class("alert alert-danger").Text("Task failed to be enqueued")
	}

	taskParametersMap := maputils.AnyToMapStringAny(taskParametersAny)

	_, err = controller.store.TaskEnqueueByAlias(task.Alias(), taskParametersMap)
	if err != nil {
		controller.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.Div().Class("alert alert-danger").Text("Task failed to be enqueued")
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Task enqueued successfully"})).
		Child(hb.Script(`setTimeout(() => {window.location.href = window.location.href;}, 3000);`))

}
