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

func queueParameters(logger slog.Logger, store taskstore.StoreInterface) *queueParametersCeontroller {
	return &queueParametersCeontroller{
		logger: logger,
		store:  store,
	}
}

type queueParametersCeontroller struct {
	logger slog.Logger
	store  taskstore.StoreInterface
}

func (c *queueParametersCeontroller) ToTag(w http.ResponseWriter, r *http.Request) hb.TagInterface {
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

	return c.modal(data)
}

func (c *queueParametersCeontroller) modal(data queueParametersCeontrollerData) *hb.Tag {
	modalID := `ModalQueueDetails`
	formID := modalID + `Form`

	fieldQueueID := form.NewField(form.FieldOptions{
		Label:    "Queue ID",
		Name:     "queue_id",
		Type:     form.FORM_FIELD_TYPE_HIDDEN,
		Value:    data.queueID,
		Required: true,
	})

	fieldParameters := form.NewField(form.FieldOptions{
		Label:    "Queued Task Parameters",
		Name:     "parameters",
		Type:     form.FORM_FIELD_TYPE_TEXTAREA,
		Value:    data.queue.Parameters(),
		Required: true,
	})

	fieldParametersSize := form.NewField(form.FieldOptions{
		Type:  form.FORM_FIELD_TYPE_RAW,
		Value: hb.Style(`#` + formID + ` textarea[name="parameters"] { height: 200px; }`).ToHTML(),
	})

	formUpdate := form.NewForm(form.FormOptions{
		ID: modalID + "Form",
		Fields: []form.FieldInterface{
			fieldParametersSize,
			fieldParameters,
			fieldQueueID,
		},
	})

	modalCloseScript := `document.getElementById('` + modalID + `').remove();document.getElementById('ModalBackdrop').remove();`
	butonModalClose := hb.Button().Type("button").
		Class("btn-close").
		Data("bs-dismiss", "modal").
		OnClick(modalCloseScript)

	buttonCancel := hb.Button().
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		HTML("Cancel").
		Class("btn btn-secondary float-start").
		OnClick(modalCloseScript)

	buttonOk := hb.Button().
		Child(hb.I().Class("bi bi-check-circle me-2")).
		HTML("OK").
		Class("btn btn-success float-end").
		OnClick(modalCloseScript)

	modal := bs.Modal().
		ID(modalID).
		Class("fade show").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						hb.Heading5().
							Text("Queue Details").
							Style(`padding: 0px; margin: 0px;`),
						butonModalClose,
					}),

					bs.ModalBody().
						Child(formUpdate.Build()),

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

func (c *queueParametersCeontroller) prepareData(r *http.Request) (data queueParametersCeontrollerData, err error) {
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

type queueParametersCeontrollerData struct {
	request *http.Request
	queueID string
	queue   taskstore.QueueInterface
}
