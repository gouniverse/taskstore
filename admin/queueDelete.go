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

func queueDelete(logger slog.Logger, store taskstore.StoreInterface) *queueDeleteCeontroller {
	return &queueDeleteCeontroller{
		logger: logger,
		store:  store,
	}
}

type queueDeleteCeontroller struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *queueDeleteCeontroller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
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

func (c *queueDeleteCeontroller) formSubmitted(data queueDeleteCeontrollerData) hb.TagInterface {
	if err := c.store.QueueSoftDelete(data.queue); err != nil {
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
			Text:              "Queue successfully deleted.",
			Position:          "top-right",
			ShowCancelButton:  false,
			ShowConfirmButton: false,
		})).
		Child(hb.Script(`setTimeout(function(){window.location.href = window.location.href}, 2000);`))
}

func (c *queueDeleteCeontroller) modal(data queueDeleteCeontrollerData) *hb.Tag {
	fieldDanger := form.NewField(form.FieldOptions{
		Type: form.FORM_FIELD_TYPE_RAW,
		Value: hb.Wrap().
			Child(hb.Paragraph().
				Child(hb.Text(`You are about to permanently delete this queued task:`))).
			Child(hb.Paragraph().
				Style("font-weight: bold;").
				Child(hb.Text(`Ref. "` + data.queue.ID() + `"`))).
			Child(hb.Paragraph().
				Child(hb.Text(`Are you sure you want to proceed?`))).
			Child(hb.Paragraph().
				Class("text-danger").
				Child(hb.Text(`Warning: This action cannot be undone.`))).
			ToHTML(),
		Required: true,
	})

	fieldQueueID := form.NewField(form.FieldOptions{
		Label:    "Queue ID",
		Name:     "queue_id",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.queueID,
		Required: true,
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: "FormQueueDelete",
		Fields: []form.FieldInterface{
			fieldDanger,
			fieldQueueID,
		},
	})

	modalCloseScript := `document.getElementById('ModalQueueDelete').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonDelete := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("Delete").
		Class("btn btn-danger float-end").
		HxInclude(`#ModalQueueDelete`).
		HxPost(url(data.request, pathQueueDelete, nil)).
		HxTarget("body").
		HxSwap("beforeend")

	modal := bs.Modal().
		ID("ModalQueueDelete").
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Delete Queue").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formUpdate.Build()),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonDelete),
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

func (c *queueDeleteCeontroller) prepareData(r *http.Request) (data queueDeleteCeontrollerData, err error) {
	data.request = r

	data.queueID = utils.Req(r, "queue_id", "")

	if data.queueID == "" {
		return data, errors.New("queue_id is required")
	}

	data.queue, err = c.store.QueueFindByID(data.queueID)

	if err != nil {
		return data, err
	}

	if data.queue == nil {
		return data, errors.New("queue not found")
	}

	return data, nil
}

type queueDeleteCeontrollerData struct {
	request *http.Request
	queueID string
	queue   taskstore.QueueInterface
}
