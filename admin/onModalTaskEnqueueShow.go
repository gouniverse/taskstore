package admin

import (
	"net/http"

	"github.com/gouniverse/hb"
)

func (controller *queueManagerController) taskEnqueueModalShow(r *http.Request) hb.TagInterface {
	return controller.modalTaskEnqueue(r)
}
