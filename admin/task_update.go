package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/utils"
)

func taskUpdate(logger slog.Logger, store taskstore.StoreInterface) *taskUpdateController {
	return &taskUpdateController{
		logger: logger,
		store:  store,
	}
}

type taskUpdateController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *taskUpdateController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              err.Error(),
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modal(data)
}

func (c *taskUpdateController) formSubmitted(data taskUpdateControllerData) hb.TagInterface {
	if data.formTitle == "" {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              "Title is required.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if data.formAlias == "" {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              "Alias is required.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	if data.formStatus == "" {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              "Status is required.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	data.task.
		SetTitle(data.formTitle).
		SetAlias(data.formAlias).
		SetStatus(data.formStatus).
		SetDescription(data.formDescription)

	err := c.store.TaskUpdate(data.task)

	if err != nil {
		return hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              err.Error(),
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Success",
			Text:              "Task successfully updated.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *taskUpdateController) modal(data taskUpdateControllerData) *hb.Tag {
	fieldTitle := form.NewField(form.FieldOptions{
		Label:    "Title",
		Name:     "title",
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.formTitle,
		Help:     "The title of the task as displayed in the dashboard.",
		Required: true,
	})

	fieldAlias := form.NewField(form.FieldOptions{
		Label:    "Alias / Command Name",
		Name:     "alias",
		Type:     form.FORM_FIELD_TYPE_STRING,
		Value:    data.formAlias,
		Help:     "The alias / the command name of the task. Should be unique.",
		Required: true,
	})

	fieldStatus := form.NewField(form.FieldOptions{
		Label:    "Status",
		Name:     "status",
		Type:     form.FORM_FIELD_TYPE_SELECT,
		Value:    data.formStatus,
		Help:     "The status of the task.",
		Required: true,
		Options: []form.FieldOption{
			{
				Value: "-- select status --",
				Key:   "",
			},
			{
				Value: "Active",
				Key:   taskstore.TaskStatusActive,
			},
			{
				Value: "Inactive",
				Key:   taskstore.TaskStatusCanceled,
			},
		},
	})

	fieldDescription := form.NewField(form.FieldOptions{
		Label:    "Description",
		Name:     "description",
		Type:     form.FORM_FIELD_TYPE_TEXTAREA,
		Value:    data.formDescription,
		Help:     "The description of the task.",
		Required: true,
	})

	fieldTaskID := form.NewField(form.FieldOptions{
		Label:    "Task ID",
		Name:     "task_id",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.taskID,
		Required: true,
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: "FormTaskUpdate",
		Fields: []form.FieldInterface{
			fieldTitle,
			fieldAlias,
			fieldStatus,
			fieldDescription,
			fieldTaskID,
		},
	})

	modalCloseScript := `document.getElementById('ModalTaskUpdate').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonUpdate := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("Save").
		Class("btn btn-success float-end").
		HxInclude(`#ModalTaskUpdate`).
		HxPost(url(data.request, pathTaskUpdate, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalTaskUpdate").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Edit Task").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formUpdate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonUpdate),
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

func (c *taskUpdateController) prepareData(r *http.Request) (data taskUpdateControllerData, err error) {
	data.request = r

	data.taskID = utils.Req(r, "task_id", "")

	if data.taskID == "" {
		return data, errors.New("task_id is required")
	}

	data.task, err = c.store.TaskFindByID(data.taskID)

	if err != nil {
		return data, err
	}

	if data.task == nil {
		return data, errors.New("task not found")
	}

	if r.Method == http.MethodGet {
		data.formAlias = data.task.Alias()
		data.formDescription = data.task.Description()
		data.formStatus = data.task.Status()
		data.formTitle = data.task.Title()
	}

	if r.Method == http.MethodPost {
		data.formAlias = utils.Req(r, "alias", "")
		data.formDescription = utils.Req(r, "description", "")
		data.formStatus = utils.Req(r, "status", "")
		data.formTitle = utils.Req(r, "title", "")
	}

	return data, nil
}

type taskUpdateControllerData struct {
	request *http.Request
	taskID  string
	task    taskstore.TaskInterface

	formAlias       string
	formDescription string
	formStatus      string
	formTitle       string
}
