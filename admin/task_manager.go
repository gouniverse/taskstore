package admin

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/utils"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

const actionModalTaskFilterShow = "modal_task_filter_show"

func taskManager(logger slog.Logger, store taskstore.StoreInterface, layout Layout) *taskManagerController {
	return &taskManagerController{
		logger: logger,
		store:  store,
		layout: layout,
	}
}

type taskManagerController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
	layout Layout
}

func (c *taskManagerController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, errorMessage := c.prepareData(r)

	c.layout.SetTitle("Task Manager | Zeppelin")

	if errorMessage != "" {
		c.layout.SetBody(hb.Div().
			Class("alert alert-danger").
			Text(errorMessage).ToHTML())

		return hb.Raw(c.layout.Render(w, r))
	}

	htmxScript := `setTimeout(() => {
		if (!window.htmx) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Htmx_2_0_0() + `';
		}
	}, 1000);`

	swalScript := `setTimeout(() => {
		if (!window.Swal) {
			let script = document.createElement('script');
			document.head.appendChild(script);
			script.type = 'text/javascript';
			script.src = '` + cdn.Sweetalert2_11() + `';
		}
	}, 1000);`

	c.layout.SetBody(c.page(data).ToHTML())
	c.layout.SetScripts([]string{htmxScript, swalScript})

	return hb.Raw(c.layout.Render(w, r))
}

func (controller *taskManagerController) page(data taskManagerControllerData) hb.TagInterface {
	adminHeader := adminHeader(controller.store, &controller.logger, data.request)
	breadcrumbs := breadcrumbs(data.request, []Breadcrumb{
		{
			Name: "Task Manager",
			URL:  url(data.request, pathTaskManager, map[string]string{}),
		},
	})

	buttonTaskCreate := hb.Button().
		Class("btn btn-primary float-end").
		Child(hb.I().Class("bi bi-plus-circle").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("New Task").
		HxGet(url(data.request, pathTaskCreate, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	title := hb.Heading1().
		HTML("Zeppelin. Task Manager").
		Child(buttonTaskCreate)

	return hb.Div().
		Class("container").
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(adminHeader).
		Child(hb.HR()).
		Child(title).
		Child(controller.tableRecords(data))
}

func (controller *taskManagerController) tableRecords(data taskManagerControllerData) hb.TagInterface {
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
						Child(controller.sortableColumnLabel(data, "Created", taskstore.COLUMN_CREATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						Child(controller.sortableColumnLabel(data, "Modified", taskstore.COLUMN_UPDATED_AT)).
						Style("width: 1px;cursor: pointer;"),
					hb.TH().
						HTML("Actions"),
				}),
			}),
			hb.Tbody().Children(lo.Map(data.recordList, func(task taskstore.TaskInterface, _ int) hb.TagInterface {
				taskName := task.Title()
				taskAlias := task.Alias()

				buttonDelete := hb.Button().
					Class("btn btn-sm btn-danger").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-trash")).
					Title("Delete task").
					HxGet(url(data.request, pathTaskDelete, map[string]string{
						"task_id": task.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				buttonUpdate := hb.Button().
					Class("btn btn-sm btn-success").
					Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
					Child(hb.I().Class("bi bi-pencil-square")).
					Title("Edit task").
					HxGet(url(data.request, pathTaskUpdate, map[string]string{
						"task_id": task.ID(),
					})).
					HxTarget("body").
					HxSwap("beforeend")

				createdAtDate := task.CreatedAtCarbon().Format("d M Y")
				createdAtTime := task.CreatedAtCarbon().ToTimeString()
				updatedAtDate := task.UpdatedAtCarbon().Format("d M Y")
				updatedAtTime := task.UpdatedAtCarbon().ToTimeString()

				status := hb.Span().
					Style(`font-weight: bold;`).
					// StyleIf(taskdTask.IsFailed(), `color:red;`).
					HTML(task.Status())

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
							Text(task.ID())),
					hb.TD().
						Child(status),
					hb.TD().
						Child(hb.Div().Text(createdAtDate)).
						Child(hb.Div().Text(createdAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Child(hb.Div().Text(updatedAtDate)).
						Child(hb.Div().Text(updatedAtTime)).
						Style("white-space: nowrap; font-size: 13px;"),
					hb.TD().
						Style("text-align: center;").
						Child(buttonUpdate).
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

func (controller *taskManagerController) sortableColumnLabel(data taskManagerControllerData, tableLabel string, columnName string) hb.TagInterface {
	isSelected := strings.EqualFold(data.sortBy, columnName)

	direction := lo.If(data.sortOrder == sb.ASC, sb.DESC).Else(sb.ASC)

	if !isSelected {
		direction = sb.ASC
	}

	link := url(data.request, pathTaskManager, map[string]string{
		"controller":     pathTaskManager,
		"page":           "0",
		"by":             columnName,
		"sort":           direction,
		"date_from":      data.formCreatedFrom,
		"date_to":        data.formCreatedTo,
		"status":         data.formStatus,
		"filter_task_id": data.formTaskID,
		"task_id":        data.formTaskID,
	})
	return hb.Hyperlink().
		HTML(tableLabel).
		Child(controller.sortingIndicator(columnName, data.sortBy, direction)).
		Href(link)
}

func (controller *taskManagerController) sortingIndicator(columnName string, sortByColumnName string, sortOrder string) hb.TagInterface {
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

func (controller *taskManagerController) tableFilter(data taskManagerControllerData) hb.TagInterface {
	buttonFilter := hb.Button().
		Class("btn btn-sm btn-info text-white me-2").
		Style("margin-bottom: 2px; margin-left:2px; margin-right:2px;").
		Child(hb.I().Class("bi bi-filter me-2")).
		Text("Filters").
		HxPost(url(data.request, pathTaskManager, map[string]string{
			"action":       actionModalTaskFilterShow,
			"name":         data.formName,
			"status":       data.formStatus,
			"task_id":      data.formTaskID,
			"created_from": data.formCreatedFrom,
			"created_to":   data.formCreatedTo,
		})).
		HxTarget("body").
		HxSwap("beforeend")

	description := []string{
		hb.Span().HTML("Showing tasks").Text(" ").ToHTML(),
	}

	if data.formStatus != "" {
		description = append(description, hb.Span().Text("with status: "+data.formStatus).ToHTML())
	} else {
		description = append(description, hb.Span().Text("with status: any").ToHTML())
	}

	if data.formName != "" {
		description = append(description, hb.Span().Text("and name: "+data.formName).ToHTML())
	}

	if data.formTaskID != "" {
		description = append(description, hb.Span().Text("and task id: "+data.formTaskID).ToHTML())
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

func (controller *taskManagerController) tablePagination(data taskManagerControllerData, count int, page int, perPage int) hb.TagInterface {
	url := url(data.request, pathTaskManager, map[string]string{
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

func (controller *taskManagerController) prepareData(r *http.Request) (data taskManagerControllerData, errorMessage string) {
	var err error
	initialPerPage := 20
	data.request = r
	data.action = utils.Req(r, "action", "")
	data.taskID = utils.Req(r, "task_id", "")

	data.page = utils.Req(r, "page", "0")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(utils.Req(r, "per_page", cast.ToString(initialPerPage)))
	data.sortOrder = utils.Req(r, "sort", sb.DESC)
	data.sortBy = utils.Req(r, "by", taskstore.COLUMN_CREATED_AT)

	data.formCreatedFrom = utils.Req(r, "filter_created_from", "")
	data.formCreatedTo = utils.Req(r, "filter_created_to", "")
	data.formName = utils.Req(r, "filter_name", "")
	data.formTaskID = utils.Req(r, "filter_task_id", "")
	data.formStatus = utils.Req(r, "filter_status", "")

	data.recordList, data.recordCount, err = controller.fetchRecordList(data)

	if err != nil {
		controller.logger.Error("At taskManagerController > prepareData", "error", err.Error())
		return data, "error retrieving web tasks"
	}

	return data, ""
}

func (controller *taskManagerController) fetchRecordList(data taskManagerControllerData) (records []taskstore.TaskInterface, recordCount int64, err error) {
	taskIDs := []string{}

	if data.formTaskID != "" {
		taskIDs = append(taskIDs, data.formTaskID)
	}

	// if data.formCreatedFrom != "" {
	// 	query.CreatedAtGte = data.formCreatedFrom + " 00:00:00"
	// }

	// if data.formCreatedTo != "" {
	// 	query.CreatedAtLte = data.formCreatedTo + " 23:59:59"
	// }

	query := taskstore.TaskQuery().
		SetLimit(data.perPage).
		SetOffset(data.pageInt * data.perPage).
		SetOrderBy(data.sortBy).
		SetSortOrder(data.sortOrder)

	if len(taskIDs) > 0 {
		query = query.SetIDIn(taskIDs)
	}

	if data.formStatus != "" {
		query = query.SetStatus(data.formStatus)
	}

	// if data.formName != "" {
	// 	query = query.SetNameLike(data.formName)
	// }

	recordList, err := controller.store.TaskList(query)

	if err != nil {
		return records, 0, err
	}

	recordCount, err = controller.store.TaskCount(query)

	if err != nil {
		return []taskstore.TaskInterface{}, 0, err
	}

	return recordList, recordCount, nil
}

type taskManagerControllerData struct {
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
	formTaskID      string

	recordList  []taskstore.TaskInterface
	recordCount int64

	taskID string
}
