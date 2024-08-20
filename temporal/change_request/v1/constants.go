package v1

const (
	ChangeRequestTaskQueueGoV1 = "change_request_task_queue_go_v1"
)

type ChangeRequestWorkflowEventsV1 string

const (
	ChangeRequestWorkflowEventV1UpdateChanges               ChangeRequestWorkflowEventsV1 = "EventUpdateChanges"
	ChangeRequestWorkflowEventV1Discard                     ChangeRequestWorkflowEventsV1 = "EventDiscard"
	ChangeRequestWorkflowEventV1SubmitForReview             ChangeRequestWorkflowEventsV1 = "EventSubmitForReview"
	ChangeRequestWorkflowEventV1IntermediateReviewCallbacks ChangeRequestWorkflowEventsV1 = "EventIntermediateReviewCallback"
	ChangeRequestWorkflowEventV1FinalReview                 ChangeRequestWorkflowEventsV1 = "EventFinalReviewCallback"
)

type ChangeRequestWorkflowQueryV1 string

const (
	ChangeRequestWorkflowQueryV1GetStatus ChangeRequestWorkflowQueryV1 = "GetStatus"
)
