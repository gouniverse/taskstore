package admin

var endpoint = "" // initialized in admin.go

const pathHome = "home"

const pathQueueCreate = "queue-create"
const pathQueueDelete = "queue-delete"
const pathQueueDetails = "queue-details"
const pathQueueManager = "queue-manager"
const pathQueueParameters = "queue-parameters"
const pathQueueRequeue = "queue-requeue"
const pathQueueTaskRestart = "queue-task-restart"

// const pathQueueUpdate = "queue-update"

const pathTaskCreate = "task-create"
const pathTaskManager = "task-manager"
const pathTaskUpdate = "task-update"
const pathTaskDelete = "task-delete"

const actionModalQueuedTaskFilterShow = "modal-queued-task-filter-show"

// const actionModalQueuedTaskRequeueShow = "modal-queued-task-requeue-show"
// const actionModalQueuedTaskRequeueSubmitted = "modal-queued-task-requeue-submitted"
const actionModalQueuedTaskRestartShow = "modal-queue-task-restart-show"
const actionModalQueuedTaskRestartSubmitted = "modal-queue-task-restart-submitted"

// const actionModalQueuedTaskEnqueueShow = "modal-queued-task-enqueue-show"
// const actionModalQueuedTaskEnqueueSubmitted = "modal-task-enqueue-submitted"
