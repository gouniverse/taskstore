package admin

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/maputils"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const actionModalQueueFilterShow = "modal_queue_filter_show"

func queueManager(logger slog.Logger, store taskstore.StoreInterface, layout Layout) *queueManagerController {
	return &queueManagerController{
		logger: logger,
		store:  store,
		layout: layout,
	}
}

type queueManagerController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
	layout Layout
}

func (c *queueManagerController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.layout.SetTitle("Queue Manager | Zepelin")

	if errorMessage != "" {
		c.layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.layout.Render(w, r))
	}

	if data.action == actionModalQueueFilterShow {
		return c.onModalRecordFilterShow(data)
	}

	// if data.action == actionModalQueuedTaskDeleteShow {
	// 	return c.onModalTaskDeleteShow(r)
	// }

	// if data.action == actionModalQueuedTaskDeleteSubmitted {
	// 	return c.onModalTaskDeleteSubmitted(r)
	// }

	if data.action == actionModalQueuedTaskEnqueueShow {
		return c.onModalTaskEnqueueShow(r)
	}

	if data.action == actionModalQueuedTaskEnqueueSubmitted {
		return c.onModalTaskEnqueueSubmitted(r)
	}

	if data.action == actionModalQueuedTaskDetailsShow {
		return c.onModalTaskDetailsShow(data.queueID)
	}

	if data.action == actionModalQueuedTaskFilterShow {
		// return c.onModalQueuedTaskFilterShow(data)
	}

	if data.action == actionModalQueuedTaskParametersShow {
		return c.onModalTaskParametersShow(data.queueID)
	}

	if data.action == actionModalQueuedTaskRequeueShow {
		return c.onModalTaskRequeueShow(r, data.queueID)
	}

	if data.action == actionModalQueuedTaskRequeueSubmitted {
		return c.onModalTaskRequeueSubmitted(r)
	}

	htmxScript := `setTimeout(() => async function() {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
			await script.onload
		}
	}, 1000);`

	swalScript := `setTimeout(() => async function() {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
			await script.onload
		}
	}, 1000);`

	c.layout.SetBody(c.page(data).ToHTML())
	c.layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.layout.Render(w, r))
}

func (c *queueManagerController) onModalTaskDetailsShow(queueID string) *hb.Tag {
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

	return c.modalTaskDetails(queue.Details())
}

func (controller *queueManagerController) onModalTaskEnqueueShow(r *http.Request) hb.TagInterface {
	return controller.modalTaskEnqueue(r)
}

func (c *queueManagerController) onModalTaskEnqueueSubmitted(r *http.Request) hb.TagInterface {
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

	task, err := c.store.TaskFindByID(taskID)

	if err != nil {
		c.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	if task == nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: "Task not found"})
	}

	taskParametersAny, err := utils.FromJSON(taskParameters, map[string]interface{}{})

	if err != nil {
		c.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	taskParametersMap := maputils.AnyToMapStringAny(taskParametersAny)

	_, err = c.store.TaskEnqueueByAlias(task.Alias(), taskParametersMap)

	if err != nil {
		c.logger.Error("At adminTasks > onModalTaskEnqueueSubmitted", "error", err.Error())
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Task enqueued successfully"})).
		Child(hb.Script(`setTimeout(() => {window.location.href = window.location.href;}, 3000);`))
}

func (c *queueManagerController) onModalTaskParametersShow(queueID string) hb.TagInterface {
	if queueID == "" {
		return hb.Div().Class("alert alert-danger").Text("queue id is required")
	}

	queue, err := c.store.QueueFindByID(queueID)

	if err != nil {
		c.logger.Error("At taskadmin > onModalQueuedTaskParametersShow", "error", err.Error())
		return hb.Div().Class("alert alert-danger").Text("Error retrieving queued task")
	}

	if queue == nil {
		return hb.Div().Class("alert alert-danger").Text("Queue task not found")
	}

	return c.modalTaskParameters(queue.Parameters())
}

func (controller *queueManagerController) onModalTaskRequeueShow(r *http.Request, queueID string) hb.TagInterface {
	queue, err := controller.store.QueueFindByID(queueID)

	if err != nil {
		controller.logger.Error("At taskadmin > onModalQueuedTaskRequeueShow", "error", err.Error())
		return hb.Div().Class("alert alert-danger").Text("Error retrieving queued task")
	}

	if queue == nil {
		return hb.Div().Class("alert alert-danger").Text("Queued task not found")
	}

	return controller.modalTaskRequeue(r, queue)
}

func (controller *queueManagerController) onModalTaskRequeueSubmitted(r *http.Request) hb.TagInterface {
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

func (*queueManagerController) onModalRecordFilterShow(data queueManagerControllerData) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Filters").
		Style(`margin:0px;padding:0px;`)

	buttonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check me-2")).
		HTML("Apply").
		Class("btn btn-primary float-end").
		OnClick(`FormFilters.submit();` + modalCloseScript)

	fieldQueueID := form.NewField(form.FieldOptions{
		Label: "Queue ID",
		Name:  "filter_queue_id",
		Type:  form.FORM_FIELD_TYPE_STRING,
		Value: data.formQueueID,
		Help:  `Find queue by reference number (ID).`,
	})

	filterForm := form.NewForm(form.FormOptions{
		ID:        "FormFilters",
		Method:    http.MethodGet,
		ActionURL: url(data.request, pathQueueManager, map[string]string{}),
		Fields: []form.FieldInterface{
			form.NewField(form.FieldOptions{
				Label: "Status",
				Name:  "filter_status",
				Type:  form.FORM_FIELD_TYPE_SELECT,
				Help:  `The status of the queue.`,
				Value: data.formStatus,
				Options: []form.FieldOption{
					{
						Value: "",
						Key:   "",
					},
					{
						Value: "Active",
						Key:   cmsstore.SITE_STATUS_ACTIVE,
					},
					{
						Value: "Inactive",
						Key:   cmsstore.SITE_STATUS_INACTIVE,
					},
					{
						Value: "Draft",
						Key:   cmsstore.SITE_STATUS_DRAFT,
					},
				},
			}),
			form.NewField(form.FieldOptions{
				Label: "Name",
				Name:  "filter_name",
				Type:  form.FORM_FIELD_TYPE_STRING,
				Value: data.formName,
				Help:  `Filter by name.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created From",
				Name:  "filter_created_from",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedFrom,
				Help:  `Filter by creation date.`,
			}),
			form.NewField(form.FieldOptions{
				Label: "Created To",
				Name:  "filter_created_to",
				Type:  form.FORM_FIELD_TYPE_DATE,
				Value: data.formCreatedTo,
				Help:  `Filter by creation date.`,
			}),
			fieldQueueID,
			form.NewField(form.FieldOptions{
				Label: "Path",
				Name:  "path",
				Type:  form.FORM_FIELD_TYPE_HIDDEN,
				Value: pathQueueManager,
				Help:  `Path to this page.`,
			}),
		},
	}).Build()

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(filterForm),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonOk),
				}),
			}),
		})

	backdrop := hb.Div().
		ID("ModalBackdrop").
		Class("modal-backdrop fade show").
		Style("display:block;")

	return hb.Wrap().Children([]hb.TagInterface{
		modal,
		backdrop,
	})

}

func (controller *queueManagerController) page(data queueManagerControllerData) hb.TagInterface {
	adminHeader := adminHeader(controller.store, &controller.logger, data.request)
	breadcrumbs := breadcrumbs(data.request, []Breadcrumb{
		{
			Name: "Queue Manager",
			URL:  url(data.request, pathQueueManager, map[string]string{}),
		},
	})

	buttonQueueNew := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Enqueue Task").
		HxGet(url(data.request, pathQueueManager, map[string]string{
			"action": actionModalQueuedTaskEnqueueShow,
			"page":   data.page,
			"by":     data.sortBy,
			"sort":   data.sortOrder,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	title := hb.Heading1().
		HTML("Zeppelin. Queue Manager").
		Child(buttonQueueNew)

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(title).
		Child(controller.tableRecords(data))
}

func (controller *queueManagerController) tableRecords(data queueManagerControllerData) hb.TagInterface {
	table := hb.Table().
		Class("table table-striped table-hover table-bordered").
		Children([]hb.TagInterface{
			hb.Thead().Children([]hb.TagInterface{
				hb.TR().Children([]hb.TagInterface{
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Name", taskstore.COLUMN_TITLE)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Alias", taskstore.COLUMN_ALIAS)).
						Text(", ").
						Child(controller.sortableColumnLabel(data, "Reference", taskstore.COLUMN_ID)).
						Style(`cursor: pointer;`),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Status", taskstore.COLUMN_STATUS)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Started", "started_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Finished", "completed_at")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Duration", "duration")).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Created", taskstore.COLUMN_CREATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.recordList, func(queuedTask taskstore.QueueInterface, _ int) hb.TagInterface {
				task, taskExists := lo.Find(data.taskList, func(t taskstore.TaskInterface) bool {
					return t.ID() == queuedTask.TaskID()
				})

				taskName := lo.IfF(taskExists, func() string { return task.Title() }).Else("Unknown")
				taskAlias := lo.IfF(taskExists, func() string { return task.Alias() }).Else("Unknown")

				buttonDelete := hb.Button().
					Class("btn btn-sm btn-danger").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete task from queue").
					HxGet(url(data.request, pathQueueDelete, map[string]string{
						"queue_id": queuedTask.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonParameters := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-list-stars")).
					Title("See queued task parameters").
					HxPost(url(data.request, pathQueueManager, map[string]string{
						"action":   actionModalQueuedTaskParametersShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonDetails := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-info-circle-fill")).
					Title("See the details of the job run").
					HxPost(url(data.request, pathQueueManager, map[string]string{
						"action":   actionModalQueuedTaskDetailsShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonRequeue := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-arrow-repeat")).
					Title("Re-add task to queue as new job").
					HxPost(url(data.request, pathQueueManager, map[string]string{
						"action":   actionModalQueuedTaskRequeueShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonRestart := hb.Button().
					Class("btn btn-sm btn-info").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-arrow-clockwise")).
					Title("Restart this job").
					HxPost(url(data.request, pathQueueManager, map[string]string{
						"action":   actionModalQueuedTaskRestartShow,
						"queue_id": queuedTask.ID(),
						"page":     data.page,
						"by":       data.sortBy,
						"sort":     data.sortOrder,
					})).
					HxTarget("body").
					HxSwap("beforeend")

				startedAtDate := lo.IfF(queuedTask.StartedAt() != "", func() string {
					return queuedTask.StartedAtCarbon().Format("d M Y")
				}).Else("-")
				startedAtTime := lo.IfF(queuedTask.StartedAt() != "", func() string {
					return queuedTask.StartedAtCarbon().ToTimeString()
				}).Else("-")
				completeddAtDate := lo.IfF(queuedTask.CompletedAt() != "", func() string {
					return queuedTask.CompletedAtCarbon().Format("d M Y")
				}).Else("-")
				completeddAtTime := lo.IfF(queuedTask.CompletedAt() != "", func() string {
					return queuedTask.CompletedAtCarbon().ToTimeString()
				}).Else("-")
				elapsedTime := lo.IfF(queuedTask.StartedAt() != "" && queuedTask.CompletedAt() != "", func() string {
					return queuedTask.CompletedAtCarbon().DiffForHumans(queuedTask.StartedAtCarbon())
				}).Else("-")
				createdAtDate := queuedTask.CreatedAtCarbon().Format("d M Y")
				createdAtTime := queuedTask.CreatedAtCarbon().ToTimeString()

				status := hb.Span().
					Style(`font-weight: bold;`).
					StyleIf(queuedTask.IsSuccess(), `color:green;`).
					StyleIf(queuedTask.IsRunning(), `color:silver;`).
					StyleIf(queuedTask.IsQueued(), `color:blue;`).
					StyleIf(queuedTask.IsFailed(), `color:red;`).
					HTML(queuedTask.Status())

				return hb.TR().Children([]hb.TagInterface{
					hb.TD().
						Child(hb.Div().Text(taskName)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Alias: ").
							Text(taskAlias)).
						Child(hb.Div().
							Style("font-size: 11px;").
							Text("Ref: ").
							Text(queuedTask.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().Text(startedAtDate)).
						Child(hb.Div().Text(startedAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Child(hb.Div().Text(completeddAtDate)).
						Child(hb.Div().Text(completeddAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Child(hb.Div().Text(elapsedTime)).
						Style("white-space: nowrap;"),
					hb.TD().
						Child(hb.Div().Text(createdAtDate)).
						Child(hb.Div().Text(createdAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Style("text-align: center;").
						Child(buttonParameters).
						Child(buttonDetails).
						Child(buttonRequeue).
						Child(buttonRestart).
						Child(buttonDelete),
				})
			})),
		})

	return hb.Wrap().Children([]hb.TagInterface{
		controller.tableFilter(data),
		table,
		controller.tablePagination(data, int(data.recordCount), data.pageInt, data.perPage),
	})
}

func (controller *queueManagerController) sortableColumnLabel(data queueManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)

	if !isSelected {
		direction = sb.ASC
	}

	link := url(data.request, pathQueueManager, map[string]string{
		"controller":     pathQueueManager,
		"page":           "0",
		"by":             columnName,
		"sort":           direction,
		"date_from":      data.formCreatedFrom,
		"date_to":        data.formCreatedTo,
		"status":         data.formStatus,
		"filter_task_id": data.formTaskID,
		"queue_id":       data.formQueueID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *queueManagerController) sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
	isSelected := strings.EqualFold(sortByColumnName, columnName)

	direction := lo.If(isSelected && sortOrder == "asc", "up").
		ElseIf(isSelected && sortOrder == "desc", "down").
		Else("none")

	sortingIndicator := hb.Span().
		Class("sorting").
		HTMLIf(direction == "up", "&#8595;").
		HTMLIf(direction == "down", "&#8593;").
		HTMLIf(direction != "down" && direction != "up", "")

	return sortingIndicator
}

func (controller *queueManagerController) tableFilter(data queueManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info text-white me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(url(data.request, pathQueueManager, map[string]string{
			"action":       actionModalQueueFilterShow,
			"name":         data.formName,
			"status":       data.formStatus,
			"queue_id":     data.formQueueID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing queues").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formQueueID != "" {
		description = append(description, hb.Span().Text("and queue id: "+data.formQueueID).ToHTML())
	}

	// 	if data.formTaskID != "" {

	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
	} else if data.formCreatedFrom != "" {
		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
	} else if data.formCreatedTo != "" {
		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
	}

	return hb.Div().
		Class("card bg-light mb-3").
		Style("").
		Children([]hb.TagInterface{
			hb.Div().Class("card-body").
				Child(buttonFilter).
				Child(hb.Span().
					HTML(strings.Join(description, " "))),
		})
}

func (controller *queueManagerController) tablePagination(data queueManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	url := url(data.request, pathQueueManager, map[string]string{
		"status":       data.formStatus,
		"name":         data.formName,
		"created_from": data.formCreatedFrom,
		"created_to":   data.formCreatedTo,
		"by":           data.sortBy,
		"order":        data.sortOrder,
	})

	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       count,
		CurrentPageNumber: page,
		PagesToShow:       5,
		PerPage:           perPage,
		URL:               url,
	})

	return hb.Div().
		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
		HTML(pagination)
}

func (controller *queueManagerController) prepareData(r *http.Request) (data queueManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = utils.Req(r, "action", "")
	data.queueID = utils.Req(r, "queue_id", "")

	data.page = utils.Req(r, "page", "0")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(utils.Req(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = utils.Req(r, "sort", sb.DESC)
	data.sortBy = utils.Req(r, "by", cmsstore.COLUMN_CREATED_AT)

	data.formCreatedFrom = utils.Req(r, "filter_created_from", "")
	data.formCreatedTo = utils.Req(r, "filter_created_to", "")
	data.formName = utils.Req(r, "filter_name", "")
	data.formQueueID = utils.Req(r, "filter_queue_id", "")
	data.formStatus = utils.Req(r, "filter_status", "")

	data.recordList, data.recordCount, err = controller.fetchRecordList(data)

	if err != nil {
		controller.logger.Error("At queueManagerController > prepareData", "error", err.Error())
		return data, "error retrieving web queues"
	}

	data.taskList, err = controller.store.TaskList(taskstore.TaskQuery().
		SetOrderBy(taskstore.COLUMN_ALIAS).
		SetSortOrder(sb.ASC).
		SetOffset(0).
		SetLimit(100))

	if err != nil {
		controller.logger.Error("At queueManagerController > prepareData", "error", err.Error())
		return data, "error retrieving tasks"
	}

	return data, ""
}

func (controller *queueManagerController) fetchRecordList(data queueManagerControllerData) (records []taskstore.QueueInterface, recordCount int64, err error) {
	queueIDs := []string{}

	if data.formQueueID != "" {
		queueIDs = append(queueIDs, data.formQueueID)
	}

	// if data.formCreatedFrom != "" {
	// 	query.CreatedAtGte = data.formCreatedFrom + " 00:00:00"
	// }

	// if data.formCreatedTo != "" {
	// 	query.CreatedAtLte = data.formCreatedTo + " 23:59:59"
	// }

	query := taskstore.QueueQuery().
		SetLimit(data.perPage).
		SetOffset(data.pageInt * data.perPage).
		SetOrderBy(data.sortBy).
		SetSortOrder(data.sortOrder)

	if len(queueIDs) > 0 {
		query = query.SetIDIn(queueIDs)
	}

	if data.formStatus != "" {
		query = query.SetStatus(data.formStatus)
	}

	// if data.formName != "" {
	// 	query = query.SetNameLike(data.formName)
	// }

	recordList, err := controller.store.QueueList(query)

	if err != nil {
		return records, 0, err
	}

	recordCount, err = controller.store.QueueCount(query)

	if err != nil {
		return []taskstore.QueueInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type queueManagerControllerData struct {
	request *http.Request
	action  string

	page      string
	pageInt   int
	perPage   int
	sortOrder string
	sortBy    string

	formStatus      string
	formName        string
	formCreatedFrom string
	formCreatedTo   string
	formQueueID     string
	formTaskID      string

	recordList  []taskstore.QueueInterface
	recordCount int64

	queueID  string
	taskList []taskstore.TaskInterface
}
