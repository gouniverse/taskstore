package admin

// import (
// 	"github.com/gouniverse/bs"
// 	"github.com/gouniverse/hb"
// )

// func (controller *queueManagerController) modalQueuedTaskDetails(details string) *hb.Tag {
// 	modalCloseScript := `document.getElementById('ModalMessage').remove();document.getElementById('ModalBackdrop').remove();`

// 	title := hb.Heading5().
// 		Text("Queued Task Details").
// 		Style(`margin:0px;padding:0px;`)

// 	buttonModalClose := hb.Button().Type("button").
// 		Class("btn-close").
// 		Data("bs-dismiss", "modal").
// 		OnClick(modalCloseScript)

// 	buttonCancel := hb.Button().
// 		Child(hb.I().Class("bi bi-chevron-left me-2")).
// 		HTML("Cancel").
// 		Class("btn btn-secondary float-start").
// 		OnClick(modalCloseScript)

// 	buttonOk := hb.Button().
// 		Child(hb.I().Class("bi bi-check me-2")).
// 		HTML("Ok").
// 		Class("btn btn-primary float-end").
// 		OnClick(modalCloseScript)

// 	groupDetails := bs.FormGroup().
// 		Child(
// 			hb.Div().
// 				HTML("Details:").
// 				Style(`font-size:18px;color:black;font-weight:bold;`),
// 		).
// 		Child(
// 			hb.TextArea().
// 				Class("form-control").
// 				Style(`height:400px;`).
// 				Name("details").
// 				HTML(details),
// 		)

// 	modal := bs.Modal().
// 		ID("ModalMessage").
// 		Class("modal-lg fade show").
// 		Style(`display:block;position:fixed;top:50%;left:50%;transform:translate(-50%,-50%);z-index:1051;`).
// 		Children([]hb.TagInterface{
// 			bs.ModalDialog().Children([]hb.TagInterface{
// 				bs.ModalContent().Children([]hb.TagInterface{
// 					bs.ModalHeader().Children([]hb.TagInterface{
// 						title,
// 						buttonModalClose,
// 					}),

// 					bs.ModalBody().
// 						Child(
// 							groupDetails,
// 						),

// 					bs.ModalFooter().
// 						Style(`display:flex;justify-content:space-between;`).
// 						Child(buttonCancel).
// 						Child(buttonOk),
// 				}),
// 			}),
// 		})

// 	backdrop := hb.Div().
// 		ID("ModalBackdrop").
// 		Class("modal-backdrop fade show").
// 		Style("display:block;")

// 	return hb.Wrap().Children([]hb.TagInterface{
// 		modal,
// 		backdrop,
// 	})
// }
