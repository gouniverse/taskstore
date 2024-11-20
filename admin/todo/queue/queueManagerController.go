package admin

// import (
// 	"net/http"
// 	"project/config"

// 	"project/internal/helpers"
// 	"project/internal/layouts"
// 	"project/internal/links"
// 	"strings"

// 	"github.com/golang-module/carbon/v2"
// 	"github.com/gouniverse/bs"
// 	"github.com/gouniverse/cdn"
// 	"github.com/gouniverse/hb"
// 	"github.com/gouniverse/router"
// 	"github.com/gouniverse/sb"
// 	"github.com/gouniverse/taskstore"
// 	"github.com/gouniverse/utils"
// 	"github.com/samber/lo"
// 	"github.com/spf13/cast"
// )

// const ActionModalUserFilterShow = "modal_user_filter_show"

// // == CONTROLLER ==============================================================

// type queueManagerController struct{}

// var _ router.HTMLControllerInterface = (*queueManagerController)(nil)

// // == CONSTRUCTOR =============================================================

// func NewQueueManagerController() *queueManagerController {
// 	return &queueManagerController{}
// }

// func (controller *queueManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
// 	data, errorMessage := controller.prepareData(r)

// 	if errorMessage != "" {
// 		return helpers.ToFlashError(w, r, errorMessage, links.NewAdminLinks().Home(map[string]string{}), 10)
// 	}

// 	if data.action == ActionModalQueuedTaskDeleteShow {
// 		return controller.onModalQueuedTaskDeleteShow(r)
// 	}

// 	if data.action == ActionModalQueuedTaskDeleteSubmitted {
// 		return controller.onModalQueuedTaskDeleteSubmitted(r)
// 	}

// 	if data.action == ActionModalQueuedTaskEnqueueShow {
// 		return controller.taskEnqueueModalShow(r)
// 	}

// 	if data.action == ActionModalQueuedTaskEnqueueSubmitted {
// 		return controller.onModaltaskEnqueueSubmitted(r)
// 	}

// 	if data.action == ActionModalQueuedTaskDetailsShow {
// 		return controller.onModalQueuedTaskDetailsShow(data.queueID)
// 	}

// 	if data.action == ActionModalQueuedTaskFilterShow {
// 		return controller.onModalQueuedTaskFilterShow(data)
// 	}

// 	if data.action == ActionModalQueuedTaskParametersShow {
// 		return controller.onModalQueuedTaskParametersShow(data.queueID)
// 	}

// 	if data.action == ActionModalQueuedTaskRequeueShow {
// 		return controller.onModalQueuedTaskRequeueShow(data.queueID)
// 	}

// 	if data.action == ActionModalQueuedTaskRequeueSubmitted {
// 		return controller.onModalQueuedTaskRequeueSubmitted(r)
// 	}

// 	// if data.action == "queue-task-restart" {
// 	// 	content := p.queueTaskRestart(w, r, data)
// 	// 	w.Write([]byte(content))
// 	// 	return
// 	// }

// 	return layouts.NewAdminLayout(r, layouts.Options{
// 		Title:   "Tasks | Task Manager",
// 		Content: controller.page(data),
// 		ScriptURLs: []string{
// 			cdn.Htmx_2_0_0(),
// 			cdn.Sweetalert2_11(),
// 			cdn.Jquery_3_7_1(),
// 		},
// 		Styles: []string{},
// 	}).ToHTML()
// }

// func (controller *queueManagerController) page(data queueManagerControllerData) hb.TagInterface {
// 	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
// 		{
// 			Name: "Home",
// 			URL:  links.NewAdminLinks().Home(map[string]string{}),
// 		},
// 		{
// 			Name: "Tasks",
// 			URL:  links.NewAdminLinks().Tasks(map[string]string{}),
// 		},
// 		{
// 			Name: "Queue Manager",
// 			URL:  links.NewAdminLinks().Tasks(map[string]string{}),
// 		},
// 	})

// 	buttonQueueNew := hb.Button().
// 		Class("btn btn-primary float-end").
// 		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
// 		HTML("Enqueue Task").
// 		HxGet(links.NewAdminLinks().Tasks(map[string]string{
// 			"action": ActionModalQueuedTaskEnqueueShow,
// 			"page":   data.page,
// 			"by":     data.sortBy,
// 			"sort":   data.sortOrder,
// 		})).
// 		HxTarget("body").
// 		HxSwap("beforeend")
// 		// HxTarget("#QueuedTasksListTable").
// 		// HxSwap("outerHTML")

// 	title := hb.Heading1().
// 		HTML("Tasks. Queue Manager").
// 		Child(buttonQueueNew)

// 	return hb.Div().
// 		Class("container").
// 		Child(breadcrumbs).
// 		Child(hb.HR()).
// 		Child(title).
// 		Child(controller.tableQueue(data))
// }

// func (controller *queueManagerController) tableQueue(data queueManagerControllerData) hb.TagInterface {
// 	buttonSelectAll := hb.Div().
// 		Class("btn-group").
// 		Child(hb.Button().
// 			Type("button").
// 			ID("ButtonSelectAll").
// 			Class("btn btn-primary").
// 			OnClick("toggleSelection('ButtonSelectAll');").
// 			Child(hb.Input().
// 				Type("checkbox").
// 				OnClick("toggleSelection('ButtonSelectAll');"))).
// 		Child(hb.Div().Class("btn-group").
// 			Child(hb.Button().Type("button").
// 				Class("btn btn-primary dropdown-toggle").
// 				Attr("data-bs-toggle", "dropdown").
// 				Attr("aria-haspopup", "true").
// 				Attr("aria-expanded", "false").
// 				Child(hb.Span().
// 					Class("caret"))).
// 			Child(hb.UL().
// 				Class("dropdown-menu").
// 				Child(hb.LI().
// 					Child(hb.A().
// 						Class("dropdown-item").
// 						OnClick("buttonSelectAll('ButtonSelectAll');").
// 						HTML("Select All"))).
// 				Child(hb.LI().
// 					Child(hb.A().
// 						Class("dropdown-item").
// 						OnClick("buttonUnselectAll('ButtonSelectAll');").
// 						HTML("Un-select All"))).
// 				Child(hb.LI().
// 					Role("separator").
// 					Class("dropdown-divider")).
// 				Child(hb.LI().
// 					Child(hb.LI().
// 						Child(hb.A().
// 							Class("dropdown-item").
// 							OnClick("buttonDeleteSelected('ButtonSelectAll');").
// 							HTML("Delete Selected"))))))

// 	linkDelete := links.NewAdminLinks().Tasks(map[string]string{
// 		"action": ActionModalQueuedTaskDeleteShow,
// 	})
// 	buttonSelectAllScript := hb.Script(`
// 			function showModalBackupfileDelete(queueIds) {
// 			    htmx.ajax('POST', '` + linkDelete + `&queue_id=' + queueIds.join(','), {
// 					target: "body",
// 					swap: "beforeend"
// 				});
// 			}
//             function buttonDeleteSelected(id) {
//                 var values = [];
//                 var queueIds = [];
//                 $('input[name=QueueId]').each(function () {
//                     if ($(this).is(":checked")) {
//                         queueIds[queueIds.length] = $(this).val();
//                     }
//                 });
//                 showModalBackupfileDelete(queueIds);
//             }
//             function buttonSelectAll(id) {
//                 $('#' + id).find('input:checkbox').prop("checked", true);
//                 $('input[name=QueueId]').prop("checked", true);
//             }
//             function buttonUnselectAll(id) {
//                 $('#' + id).find('input:checkbox').prop("checked", false);
//                 $('input[name=QueueId]').prop("checked", false);
//             }
//             function toggleSelection(id) {
//                 var isChecked = $('#' + id).find('input:checkbox').is(":checked");
//                 if (isChecked == true) {
//                     buttonUnselectAll(id);
//                 } else {
//                     buttonSelectAll(id);
//                 }
//             }
// 	`)

// 	table := hb.Table().
// 		Class("table table-striped table-hover table-bordered").
// 		Children([]hb.TagInterface{
// 			hb.Thead().Children([]hb.TagInterface{
// 				hb.TR().
// 					Child(hb.TH(). // column for "checkbox"
// 							Style("width: 1px;").
// 							Text("")).
// 					Children([]hb.TagInterface{
// 						hb.TH().
// 							Child(controller.sortableColumnLabel(data, "Name", "name")).
// 							Text(", ").
// 							Child(controller.sortableColumnLabel(data, "Alias", "alias")).
// 							Text(", ").
// 							Child(controller.sortableColumnLabel(data, "Reference", "id")).
// 							Style(`cursor: pointer;`),
// 						hb.TH().
// 							Child(controller.sortableColumnLabel(data, "Status", "status")).
// 							Style("width: 200px;cursor: pointer;"),
// 						hb.TH().
// 							Child(controller.sortableColumnLabel(data, "Started", "started_at")).
// 							Style("width: 1px;cursor: pointer;"),
// 						hb.TH().
// 							Child(controller.sortableColumnLabel(data, "Finished", "completed_at")).
// 							Style("width: 1px;cursor: pointer;"),
// 						hb.TH().
// 							Child(controller.sortableColumnLabel(data, "Duration", "duration")).
// 							Style("width: 1px;cursor: pointer;"),

// 						hb.TH().
// 							Child(controller.sortableColumnLabel(data, "Created", "created_at")).
// 							Style("width: 1px;cursor: pointer;"),
// 						hb.TH().
// 							HTML("Actions"),
// 					}),
// 			}),
// 			hb.Tbody().Children(lo.Map(data.queuedTaskList, func(queuedTask taskstore.Queue, _ int) hb.TagInterface {
// 				task, taskExists := lo.Find(data.taskList, func(t taskstore.Task) bool {
// 					return t.ID == queuedTask.TaskID
// 				})

// 				taskName := lo.IfF(taskExists, func() string { return task.Title }).Else("Unknown")

// 				buttonDelete := hb.Button().
// 					Class("btn btn-sm btn-danger").
// 					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 					Child(hb.I().Class("bi bi-trash")).
// 					Title("Delete task from queue").
// 					HxPost(links.NewAdminLinks().Tasks(map[string]string{
// 						"action":   ActionModalQueuedTaskDeleteShow,
// 						"queue_id": queuedTask.ID,
// 						"page":     data.page,
// 						"by":       data.sortBy,
// 						"sort":     data.sortOrder,
// 					})).
// 					HxTarget("body").
// 					HxSwap("beforeend")

// 				buttonParameters := hb.Button().
// 					Class("btn btn-sm btn-info").
// 					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 					Child(hb.I().Class("bi bi-list-stars")).
// 					Title("See queued task parameters").
// 					HxPost(links.NewAdminLinks().Tasks(map[string]string{
// 						"action":   ActionModalQueuedTaskParametersShow,
// 						"queue_id": queuedTask.ID,
// 						"page":     data.page,
// 						"by":       data.sortBy,
// 						"sort":     data.sortOrder,
// 					})).
// 					HxTarget("body").
// 					HxSwap("beforeend")

// 				buttonDetails := hb.Button().
// 					Class("btn btn-sm btn-info").
// 					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 					Child(hb.I().Class("bi bi-info-circle-fill")).
// 					Title("See the details of the job run").
// 					HxPost(links.NewAdminLinks().Tasks(map[string]string{
// 						"action":   ActionModalQueuedTaskDetailsShow,
// 						"queue_id": queuedTask.ID,
// 						"page":     data.page,
// 						"by":       data.sortBy,
// 						"sort":     data.sortOrder,
// 					})).
// 					// HxTarget("#QueuedTasksListTable").
// 					// HxSelectOob("#ModalMessage").
// 					HxTarget("body").
// 					HxSwap("beforeend")

// 				buttonRequeue := hb.Button().
// 					Class("btn btn-sm btn-info").
// 					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 					Child(hb.I().Class("bi bi-arrow-repeat")).
// 					Title("Re-add task to queue as new job").
// 					HxPost(links.NewAdminLinks().Tasks(map[string]string{
// 						"action":   ActionModalQueuedTaskRequeueShow,
// 						"queue_id": queuedTask.ID,
// 						"page":     data.page,
// 						"by":       data.sortBy,
// 						"sort":     data.sortOrder,
// 					})).
// 					HxTarget("body").
// 					HxSwap("beforeend")

// 				buttonRestart := hb.Button().
// 					Class("btn btn-sm btn-info").
// 					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 					Child(hb.I().Class("bi bi-arrow-clockwise")).
// 					Title("Restart this job").
// 					HxPost(links.NewAdminLinks().Tasks(map[string]string{
// 						"action":   ActionModalQueuedTaskRestartShow,
// 						"queue_id": queuedTask.ID,
// 						"page":     data.page,
// 						"by":       data.sortBy,
// 						"sort":     data.sortOrder,
// 					})).
// 					HxTarget("body").
// 					HxSwap("beforeend")

// 				startedAtDate := lo.IfF(queuedTask.StartedAt != nil, func() string {
// 					return carbon.CreateFromStdTime(*queuedTask.StartedAt).Format("d M Y")
// 				}).Else("-")
// 				startedAtTime := lo.IfF(queuedTask.StartedAt != nil, func() string {
// 					return carbon.CreateFromStdTime(*queuedTask.StartedAt).ToTimeString()
// 				}).Else("-")
// 				completeddAtDate := lo.IfF(queuedTask.CompletedAt != nil, func() string {
// 					return carbon.CreateFromStdTime(*queuedTask.CompletedAt).Format("d M Y")
// 				}).Else("-")
// 				completeddAtTime := lo.IfF(queuedTask.CompletedAt != nil, func() string {
// 					return carbon.CreateFromStdTime(*queuedTask.CompletedAt).ToTimeString()
// 				}).Else("-")
// 				elapsedTime := lo.IfF(queuedTask.StartedAt != nil && queuedTask.CompletedAt != nil, func() string {
// 					return queuedTask.CompletedAt.Sub(*queuedTask.StartedAt).String()
// 				}).Else("-")
// 				createdAtDate := carbon.CreateFromStdTime(queuedTask.CreatedAt).Format("d M Y")
// 				createdAtTime := carbon.CreateFromStdTime(queuedTask.CreatedAt).ToTimeString()

// 				status := hb.Span().
// 					Style(`font-weight: bold;`).
// 					StyleIf(queuedTask.Status == taskstore.QueueStatusSuccess, `color:green;`).
// 					StyleIf(queuedTask.Status == taskstore.QueueStatusRunning, `color:silver;`).
// 					StyleIf(queuedTask.Status == taskstore.QueueStatusQueued, `color:blue;`).
// 					StyleIf(queuedTask.Status == taskstore.QueueStatusFailed, `color:red;`).
// 					HTML(queuedTask.Status)

// 				checkbox := hb.Input().
// 					Type("checkbox").
// 					Name("QueueId").
// 					Value(queuedTask.ID).
// 					Class("form-check-input").
// 					Style("margin-right: 5px;margin-left: 5px;")

// 				return hb.TR().
// 					Child(hb.TD().
// 						Style("padding: 0px;").
// 						Child(checkbox)).
// 					Children([]hb.TagInterface{
// 						hb.TD().
// 							Child(hb.Div().Text(taskName)).
// 							Child(hb.Div().
// 								Style("font-size: 11px;").
// 								Text("Alias: ").
// 								Text(task.Alias)).
// 							Child(hb.Div().
// 								Style("font-size: 11px;").
// 								Text("Ref: ").
// 								Text(queuedTask.ID)),
// 						hb.TD().
// 							Child(status),
// 						hb.TD().
// 							Child(hb.Div().Text(startedAtDate)).
// 							Child(hb.Div().Text(startedAtTime)).
// 							Style("white-space: nowrap; font-size: 13px;"),
// 						hb.TD().
// 							Child(hb.Div().Text(completeddAtDate)).
// 							Child(hb.Div().Text(completeddAtTime)).
// 							Style("white-space: nowrap; font-size: 13px;"),
// 						hb.TD().
// 							Child(hb.Div().Text(elapsedTime)).
// 							Style("white-space: nowrap;"),
// 						hb.TD().
// 							Child(hb.Div().Text(createdAtDate)).
// 							Child(hb.Div().Text(createdAtTime)).
// 							Style("white-space: nowrap; font-size: 13px;"),
// 						hb.TD().
// 							Style("text-align: center;").
// 							Child(buttonParameters).
// 							Child(buttonDetails).
// 							Child(buttonRequeue).
// 							Child(buttonRestart).
// 							Child(buttonDelete),
// 					})
// 			})),
// 		})

// 	return hb.Wrap().Children([]hb.TagInterface{
// 		controller.tableFilter(data),
// 		hb.Div().Child(buttonSelectAll).Style("margin-bottom: 10px;"),
// 		buttonSelectAllScript,
// 		table,
// 		controller.tablePagination(data, int(data.queuedTaskCount), data.pageInt, data.perPage),
// 	})
// }

// func (controller *queueManagerController) sortableColumnLabel(data queueManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
// 	isSelected := strings.EqualFold(data.sortBy, columnName)

// 	direction := lo.If(data.sortOrder == "asc", "desc").Else("asc")

// 	if !isSelected {
// 		direction = "asc"
// 	}

// 	link := links.NewAdminLinks().Tasks(map[string]string{
// 		"page":                "0",
// 		"by":                  columnName,
// 		"sort":                direction,
// 		"filter_created_from": data.formCreatedFrom,
// 		"filter_created_to":   data.formCreatedTo,
// 		"filter_status":       data.formStatus,
// 		"filter_task_id":      data.formTaskID,
// 		"filter_queue_id":     data.formQueueID,
// 	})
// 	return hb.Hyperlink().
// 		HTML(tableLabel).
// 		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
// 		Href(link)
// }

// func (controller *queueManagerController) sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
// 	isSelected := strings.EqualFold(sortByColumnName, columnName)

// 	direction := lo.If(isSelected && sortOrder == "asc", "up").
// 		ElseIf(isSelected && sortOrder == "desc", "down").
// 		Else("none")

// 	sortingIndicator := hb.Span().
// 		Class("sorting").
// 		HTMLIf(direction == "up", "&#8595;").
// 		HTMLIf(direction == "down", "&#8593;").
// 		HTMLIf(direction != "down" && direction != "up", "")

// 	return sortingIndicator
// }

// func (controller *queueManagerController) tableFilter(data queueManagerControllerData) hb.TagInterface {
// 	buttonFilter := hb.Button().
// 		Class("btn btn-sm btn-info me-2").
// 		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
// 		Child(hb.I().Class("bi bi-filter me-2")).
// 		Text("Filters").
// 		HxPost(links.NewAdminLinks().Tasks(map[string]string{
// 			"action":              ActionModalQueuedTaskFilterShow,
// 			"filter_task_id":      data.formTaskID,
// 			"filter_queue_id":     data.formQueueID,
// 			"filter_status":       data.formStatus,
// 			"filter_created_from": data.formCreatedFrom,
// 			"filter_created_to":   data.formCreatedTo,
// 		})).
// 		HxTarget("body").
// 		HxSwap("beforeend")

// 	description := []string{
// 		hb.Span().HTML("Showing queued tasks").Text(" ").ToHTML(),
// 	}

// 	if data.formStatus != "" {
// 		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
// 	} else {
// 		description = append(description, hb.Span().Text("with status: any").ToHTML())
// 	}

// 	if data.formQueueID != "" {
// 		description = append(description, hb.Span().Text("and queue ID: "+data.formQueueID).ToHTML())
// 	}

// 	if data.formTaskID != "" {
// 		task := lo.FindOrElse(data.taskList, taskstore.Task{}, func(task taskstore.Task) bool {
// 			return task.ID == data.formTaskID
// 		})
// 		taskTitle := lo.Ternary(lo.IsEmpty(task.Title), data.formTaskID, task.Title)

// 		description = append(description, hb.Span().Text("and task: "+taskTitle).ToHTML())
// 	}

// 	if data.formCreatedFrom != "" && data.formCreatedTo != "" {
// 		description = append(description, hb.Span().Text("and created between: "+data.formCreatedFrom+" and "+data.formCreatedTo).ToHTML())
// 	} else if data.formCreatedFrom != "" {
// 		description = append(description, hb.Span().Text("and created after: "+data.formCreatedFrom).ToHTML())
// 	} else if data.formCreatedTo != "" {
// 		description = append(description, hb.Span().Text("and created before: "+data.formCreatedTo).ToHTML())
// 	}

// 	return hb.Div().
// 		Class("card bg-light mb-3").
// 		Style("").
// 		Children([]hb.TagInterface{
// 			hb.Div().Class("card-body").
// 				Child(buttonFilter).
// 				Child(hb.Span().
// 					HTML(strings.Join(description, " "))),
// 		})
// }

// func (controller *queueManagerController) tablePagination(data queueManagerControllerData, count int, page int, perPage int) hb.TagInterface {
// 	url := links.NewAdminLinks().UsersUserManager(map[string]string{
// 		"filter_status":       data.formStatus,
// 		"filter_task_id":      data.formTaskID,
// 		"filter_queue_id":     data.formQueueID,
// 		"filter_created_from": data.formCreatedFrom,
// 		"filter_created_to":   data.formCreatedTo,
// 		"by":                  data.sortBy,
// 		"order":               data.sortOrder,
// 	})

// 	url = lo.Ternary(strings.Contains(url, "?"), url+"&page=", url+"?page=") // page must be last

// 	pagination := bs.Pagination(bs.PaginationOptions{
// 		NumberItems:       count,
// 		CurrentPageNumber: page,
// 		PagesToShow:       5,
// 		PerPage:           perPage,
// 		URL:               url,
// 	})

// 	return hb.Div().
// 		Class(`d-flex justify-content-left mt-5 pagination-primary-soft rounded mb-0`).
// 		HTML(pagination)
// }

// func (controller *queueManagerController) prepareData(r *http.Request) (data queueManagerControllerData, errorMessage string) {
// 	var err error
// 	data.request = r
// 	data.action = utils.Req(r, "action", "")
// 	data.page = utils.Req(r, "page", "0")
// 	data.pageInt = cast.ToInt(data.page)
// 	data.perPage = cast.ToInt(utils.Req(r, "per_page", "10"))
// 	data.sortOrder = utils.Req(r, "sort_order", sb.DESC)
// 	data.sortBy = utils.Req(r, "by", taskstore.COLUMN_CREATED_AT)
// 	data.formTaskID = utils.Req(r, "filter_task_id", "")
// 	data.formQueueID = utils.Req(r, "filter_queue_id", "")
// 	data.formStatus = utils.Req(r, "filter_status", "")
// 	data.formCreatedFrom = utils.Req(r, "filter_created_from", "")
// 	data.formCreatedTo = utils.Req(r, "filter_created_to", "")
// 	data.queueID = utils.Req(r, "queue_id", "")

// 	if !lo.Contains([]string{sb.ASC, sb.DESC}, data.sortOrder) {
// 		data.sortOrder = sb.DESC
// 	}

// 	if !lo.Contains([]string{
// 		taskstore.COLUMN_STARTED_AT,
// 		taskstore.COLUMN_COMPLETED_AT,
// 		taskstore.COLUMN_ID,
// 		taskstore.COLUMN_TASK_ID,
// 		taskstore.COLUMN_STATUS,
// 	}, data.sortBy) {
// 		data.sortBy = taskstore.COLUMN_CREATED_AT
// 	}

// 	taskList, queuedTaskList, queuedTaskCount, err := controller.fetchQueueList(data)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 		return data, "error retrieving users"
// 	}

// 	data.taskList = taskList
// 	data.queuedTaskList = queuedTaskList
// 	data.queuedTaskCount = queuedTaskCount

// 	return data, ""
// }

// func (controller *queueManagerController) fetchQueueList(data queueManagerControllerData) ([]taskstore.Task, []taskstore.Queue, int64, error) {
// 	// userIDs := []string{}

// 	// if data.formFirstName != "" {
// 	// 	firstNameUserIDs, err := config.BlindIndexStoreFirstName.Search(data.formFirstName, blindindexstore.SEARCH_TYPE_CONTAINS)

// 	// 	if err != nil {
// 	// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 	// 		return []userstore.User{}, 0, err
// 	// 	}

// 	// 	if len(firstNameUserIDs) == 0 {
// 	// 		return []userstore.User{}, 0, nil
// 	// 	}

// 	// 	userIDs = append(userIDs, firstNameUserIDs...)
// 	// }

// 	// if data.formLastName != "" {
// 	// 	lastNameUserIDs, err := config.BlindIndexStoreLastName.Search(data.formLastName, blindindexstore.SEARCH_TYPE_CONTAINS)

// 	// 	if err != nil {
// 	// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 	// 		return []userstore.User{}, 0, err
// 	// 	}

// 	// 	if len(lastNameUserIDs) == 0 {
// 	// 		return []userstore.User{}, 0, nil
// 	// 	}

// 	// 	userIDs = append(userIDs, lastNameUserIDs...)
// 	// }

// 	// if data.formEmail != "" {
// 	// 	emailUserIDs, err := config.BlindIndexStoreEmail.Search(data.formEmail, blindindexstore.SEARCH_TYPE_CONTAINS)

// 	// 	if err != nil {
// 	// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 	// 		return []userstore.User{}, 0, err
// 	// 	}

// 	// 	if len(emailUserIDs) == 0 {
// 	// 		return []userstore.User{}, 0, nil
// 	// 	}

// 	// 	userIDs = append(userIDs, emailUserIDs...)
// 	// }

// 	taskList, err := config.TaskStore.TaskList(taskstore.TaskQueryOptions{})

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 		return []taskstore.Task{}, []taskstore.Queue{}, 0, err
// 	}

// 	query := taskstore.QueueQueryOptions{
// 		// 	IDIn:      userIDs,
// 		// 	Offset:    data.pageInt * data.perPage,
// 		// 	Limit:     data.perPage,
// 		TaskID: data.formTaskID,
// 		Status: data.formStatus,
// 		// 	SortOrder: data.sortOrder,
// 		// 	OrderBy:   data.sortBy,
// 	}

// 	// if data.formCreatedFrom != "" {
// 	// 	query.CreatedAtGte = data.formCreatedFrom + " 00:00:00"
// 	// }

// 	// if data.formCreatedTo != "" {
// 	// 	query.CreatedAtLte = data.formCreatedTo + " 23:59:59"
// 	// }

// 	queuedTaskList, err := config.TaskStore.QueueList(query)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 		return []taskstore.Task{}, []taskstore.Queue{}, 0, err
// 	}

// 	queuedTaskCount, err := config.TaskStore.QueueCount(query)

// 	if err != nil {
// 		config.LogStore.ErrorWithContext("At queueManagerController > prepareData", err.Error())
// 		return []taskstore.Task{}, []taskstore.Queue{}, 0, err
// 	}

// 	return taskList, queuedTaskList, queuedTaskCount, nil
// }

// type queueManagerControllerData struct {
// 	request         *http.Request
// 	action          string
// 	page            string
// 	pageInt         int
// 	perPage         int
// 	sortOrder       string
// 	sortBy          string
// 	formStatus      string
// 	formTaskID      string
// 	formCreatedFrom string
// 	formCreatedTo   string
// 	formQueueID     string
// 	queueID         string
// 	taskList        []taskstore.Task
// 	queuedTaskList  []taskstore.Queue
// 	queuedTaskCount int64
// }
