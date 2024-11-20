package admin

var endpoint = "" // initialized in admin.go

const pathQueueCreate = "queue-create"
const pathQueueManager = "queue-manager"
const pathQueueUpdate = "queue-update"
const pathQueueDelete = "queue-delete"

const PathHome = "home"
const actionModalQueuedTaskDeleteShow = "modal-queued-task-delete-show"
const actionModalQueuedTaskDeleteSubmitted = "modal-queued-task-delete-submitted"
const actionModalQueuedTaskDetailsShow = "modal-queued-task-details-show"
const actionModalQueuedTaskFilterShow = "modal-queued-task-filter-show"
const actionModalQueuedTaskParametersShow = "modal-queued-task-parameters-show"
const actionModalQueuedTaskRequeueShow = "modal-queued-task-requeue-show"
const actionModalQueuedTaskRequeueSubmitted = "modal-queued-task-requeue-submitted"
const actionModalQueuedTaskRestartShow = "modal-queue-task-restart-show"
const actionModalQueuedTaskRestartSubmitted = "modal-queue-task-restart-submitted"
const actionModalQueuedTaskEnqueueShow = "modal-queued-task-enqueue-show"
const actionModalQueuedTaskEnqueueSubmitted = "modal-task-enqueue-submitted"
