package taskstore

import (
	"strings"

	"github.com/mingrammer/cfmt"
)

type BaseTaskHandler struct {
	QueuedTask *Queue // dynamic

	ErrorMessage   string
	InfoMessage    string
	SuccessMessage string
}

func (t *BaseTaskHandler) HasQueuedTask() bool {
	return t.QueuedTask != nil
}

func (t *BaseTaskHandler) LogError(message string) {
	t.ErrorMessage = message
	if t.HasQueuedTask() {
		t.QueuedTask.AppendDetails(message)
	} else {
		cfmt.Errorln(message)
	}
}

func (t *BaseTaskHandler) LogInfo(message string) {
	t.InfoMessage = message
	if t.HasQueuedTask() {
		t.QueuedTask.AppendDetails(message)
	} else {
		cfmt.Infoln(message)
	}
}

func (task *BaseTaskHandler) LogSuccess(message string) {
	task.SuccessMessage = message
	if task.HasQueuedTask() {
		task.QueuedTask.AppendDetails(message)
	} else {
		cfmt.Successln(message)
	}
}

func (t *BaseTaskHandler) GetParam(paramName string, opts TaskHandlerOptions) string {
	if opts.QueuedTask != nil {
		parameters, parametersErr := opts.QueuedTask.GetParameters()
		if parametersErr != nil {
			opts.QueuedTask.AppendDetails("Parameters JSON incorrect. " + parametersErr.Error())
			return ""
		}

		parameter, parameterExists := parameters[paramName]
		if !parameterExists {
			return ""
		}

		return parameter.(string)
	} else {
		return opts.Arguments[paramName]
	}
}

func (t *BaseTaskHandler) GetParamArray(paramName string, opts TaskHandlerOptions) []string {
	if opts.QueuedTask != nil {
		parameters, parametersErr := opts.QueuedTask.GetParameters()
		if parametersErr != nil {
			opts.QueuedTask.AppendDetails("Parameters JSON incorrect. " + parametersErr.Error())
			return []string{}
		}

		paramValues, paramExists := parameters[paramName]
		if !paramExists {
			return []string{}
		}

		paramValuesInterface := paramValues.([]interface{})
		paramValuesString := []string{}
		for _, v := range paramValuesInterface {
			paramValuesString = append(paramValuesString, v.(string))
		}
		return paramValuesString
	} else {
		return strings.Split(opts.Arguments[paramName], ";")
	}
}
