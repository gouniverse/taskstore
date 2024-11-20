package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/utils"
)

func taskCreate(logger slog.Logger, store taskstore.StoreInterface) *taskCreateController {
	return &taskCreateController{
		logger: logger,
		store:  store,
	}
}

type taskCreateController struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *taskCreateController) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
	data, err := c.prepareData(r)

	if err != nil {
		return hb.Swal(hb.SwalOptions{Title: "Error", Text: err.Error()})
	}

	if r.Method == http.MethodPost {
		return c.formSubmitted(data)
	}

	return c.modalTaskCreate(data)
}

func (c *taskCreateController) formSubmitted(data taskCreateControllerData) hb.TagInterface {
	if data.formTitle == "" {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Title is required.", Position: "top-right"})
	}

	if data.formAlias == "" {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Alias is required.", Position: "top-right"})
	}

	if data.formStatus == "" {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: "Status is required.", Position: "top-right"})
	}

	task := taskstore.NewTask().
		SetTitle(data.formTitle).
		SetAlias(data.formAlias).
		SetStatus(data.formStatus).
		SetDescription(data.formDescription)

	err := c.store.TaskCreate(task)

	if err != nil {
		return hb.Swal(hb.SwalOptions{Icon: "error", Title: "Error", Text: err.Error(), Position: "top-right"})
	}

	return hb.Wrap().
		Child(hb.Swal(hb.SwalOptions{Icon: "success", Title: "Success", Text: "Task successfully created.", Position: "top-right"})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *taskCreateController) modalTaskCreate(data taskCreateControllerData) *hb.Tag {
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

	formCreate := form.NewForm(form.FormOptions{
		ID: "ModalTaskCreate",
		Fields: []form.FieldInterface{
			fieldTitle,
			fieldAlias,
			fieldStatus,
		},
	})

	modalCloseScript := `document.getElementById('ModalTaskCreate').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonCreate := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("Create").
		Class("btn btn-success float-end").
		HxInclude(`#ModalTaskCreate`).
		HxPost(url(data.request, pathTaskCreate, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalTaskCreate").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("New Task").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formCreate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonCreate),
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

func (c *taskCreateController) prepareData(r *http.Request) (data taskCreateControllerData, err error) {
	data.request = r
	data.formAlias = utils.Req(r, "alias", "")
	data.formDescription = utils.Req(r, "description", "")
	data.formStatus = utils.Req(r, "status", "")
	data.formTitle = utils.Req(r, "title", "")

	return data, nil
}

type taskCreateControllerData struct {
	request *http.Request

	formAlias       string
	formDescription string
	formStatus      string
	formTitle       string
}
