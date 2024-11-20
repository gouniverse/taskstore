package admin

import (
	"log/slog"
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/taskstore"
	"github.com/spf13/cast"
)

func adminHeader(store taskstore.StoreInterface, logger *slog.Logger, r *http.Request) hb.TagInterface {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(url(r, pathQueueManager, nil)).
		Class("nav-link")
	linkQueue := hb.Hyperlink().
		HTML("Queue").
		Href(url(r, pathQueueManager, nil)).
		Class("nav-link")
	linkTasks := hb.Hyperlink().
		HTML("Tasks").
		Href(url(r, pathQueueManager, nil)).
		Class("nav-link")

	queueCount, err := store.QueueCount(taskstore.QueueQuery())

	if err != nil {
		logger.Error(err.Error())
		queueCount = -1
	}

	taskCount, err := store.TaskCount(taskstore.TaskQuery())

	if err != nil {
		logger.Error(err.Error())
		taskCount = -1
	}

	ulNav := hb.NewUL().Class("nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkHome))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkQueue.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(queueCount)))))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkTasks.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(taskCount)))))

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
}
