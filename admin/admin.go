package admin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/utils"
)

type UIOptions struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Logger         *slog.Logger
	Store          taskstore.StoreInterface
}

func UI(options UIOptions) (hb.TagInterface, error) {
	if options.ResponseWriter == nil {
		return nil, errors.New("options.ResponseWriter is required")
	}

	if options.Request == nil {
		return nil, errors.New("options.Request is required")
	}

	if options.Store == nil {
		return nil, errors.New("options.Store is required")
	}

	if options.Logger == nil {
		return nil, errors.New("options.Logger is required")
	}

	admin := &admin{
		response: options.ResponseWriter,
		request:  options.Request,
		store:    options.Store,
		logger:   *options.Logger,
	}
	return admin.handler(), nil
}

type admin struct {
	response http.ResponseWriter
	request  *http.Request
	store    taskstore.StoreInterface
	logger   slog.Logger
}

func (a *admin) handler() hb.TagInterface {
	controller := utils.Req(a.request, "controller", "")

	if controller == "" {
		controller = pathQueueManager
	}

	if controller == pathTaskCreate {
		return taskCreate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskManager {
		return taskManager(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathTaskUpdate {
		return taskUpdate(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueCreate {
		// 	return queueCreateUi(a.logger, a.store).ToTag(a.response, a.request)
		return hb.Div().Child(hb.H1().HTML(controller))
	}

	if controller == pathQueueDelete {
		// 	return queueDeleteUi(a.logger, a.store).ToTag(a.response, a.request)
		return hb.Div().Child(hb.H1().HTML(controller))
	}

	if controller == pathQueueManager {
		return queueManagerUi(a.logger, a.store).ToTag(a.response, a.request)
	}

	if controller == pathQueueUpdate {
		// 	return queueUpdateUi(a.logger, a.store).ToTag(a.response, a.request)
		return hb.Div().Child(hb.H1().HTML(controller))
	}

	return hb.Div().Child(hb.H1().HTML(controller))
	// redirect(a.response, a.request, url(a.request, pathQueueManager, map[string]string{}))
	// return nil
}
