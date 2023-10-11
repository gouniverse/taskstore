package taskstore

// type TaskHandlerOptions struct {
// 	QueuedTask *Queue
// 	Arguments  map[string]string
// }

// func (opts TaskOptions) HasTask() bool {
// 	return opts.QueuedTask != nil
// }

// func (opts TaskOptions) LogError(message string) {
// 	if opts.HasTask() {
// 		opts.QueuedTask.AppendDetails(message)
// 		config.TaskStore.QueueUpdate(opts.QueuedTask)
// 	} else {
// 		cfmt.Errorln(message)
// 	}
// }

// func (opts TaskOptions) LogInfo(message string) {
// 	if opts.HasTask() {
// 		opts.QueuedTask.AppendDetails(message)
// 		config.TaskStore.QueueUpdate(opts.QueuedTask)
// 	} else {
// 		cfmt.Infoln(message)
// 	}
// }

// func (opts TaskOptions) LogSuccess(message string) {
// 	if opts.HasTask() {
// 		opts.QueuedTask.AppendDetails(message)
// 		config.TaskStore.QueueUpdate(opts.QueuedTask)
// 	} else {
// 		cfmt.Successln(message)
// 	}
// }
