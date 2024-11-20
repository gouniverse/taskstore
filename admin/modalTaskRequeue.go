package admin

import (
	"net/http"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/taskstore"
)

func (controller *queueManagerController) modalTaskRequeue(r *http.Request, queuedTask taskstore.QueueInterface) *hb.Tag {
	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

	title := hb.Heading5().
		Text("Queued Task Requeue").
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

	buttonRequeue := hb.Button().
		Child(hb.I().Class("bi bi-arrow-clockwise me-2")).
		HTML("Requeue").
		Class("btn btn-primary float-end").
		HxPost(url(r, pathQueueManager, map[string]string{
			"action":  actionModalQueuedTaskRequeueSubmitted,
			"task_id": queuedTask.TaskID(),
		})).
		HxInclude("#ModalMessage").
		HxTarget("body").
		HxSwap("beforeend")

	divInfo := hb.Div().
		Class("alert alert-info").
		Text(`A new task will be created with the following parameters. You may  edit the parameters if necessary`)

	groupParameters := bs.FormGroup().
		Child(
			hb.Div().
				HTML("Requeue Parameters:").
				Style(`font-size:18px;color:black;font-weight:bold;`),
		).
		Child(
			hb.TextArea().
				Class("form-control").
				Style(`height:300px;`).
				Name("task_parameters").
				HTML(queuedTask.Parameters()),
		).
		Child(
			hb.Div().
				Class("form-text text-muted mb-3").
				Text(`Must be valid JSON.`),
		)

	modal := bs.Modal().
		ID("ModalMessage").
		Class("fade show modal-lg").
		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
		Children([]hb.TagInterface{
			bs.ModalDialog().Children([]hb.TagInterface{
				bs.ModalContent().Children([]hb.TagInterface{
					bs.ModalHeader().Children([]hb.TagInterface{
						title,
						buttonModalClose,
					}),

					bs.ModalBody().
						Child(divInfo).
						Child(
							groupParameters,
						),

					bs.ModalFooter().
						Style(`display:flex;justify-content:space-between;`).
						Child(buttonCancel).
						Child(buttonRequeue),
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
