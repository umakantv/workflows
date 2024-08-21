package workflow

import (
	"log"
	"time"

	"github.com/umakantv/workflows/common/model"
	v1 "github.com/umakantv/workflows/temporal/change_request/v1"
	"go.temporal.io/sdk/workflow"
)

// ChangeRequestWorkflow is a Temporal workflow that orchestrates the change request process.
func ChangeRequestWorkflowV1(ctx workflow.Context, id string) error {
	// @@@SNIPBEGIN temporal/workflow/v1/workflow.go ChangeRequestWorkflow

	// Setup the workflow state
	workflow.GetLogger(ctx).Info("Starting change request workflow", "request", id)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	// Signals we need to handle
	// 1. Update Changes (n signals)
	// 2. Discard
	// 3. Submit for Review
	// 4. Intermediate Approval Callbacks
	// 5. Final Review - Approval/Rejection

	// Start the workflow

	discarded := false
	submitted := false

	for {
		discardSubmitOrAddChangesSelector := workflow.NewSelector(ctx)

		discardCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1Discard))
		discardSubmitOrAddChangesSelector.AddReceive(discardCh, func(ch workflow.ReceiveChannel, more bool) {
			ch.Receive(ctx, nil)
			workflow.GetLogger(ctx).Info("Discarding change request", "request", id)
			err := workflow.ExecuteActivity(ctx, string(v1.ChangeRequestWorkflowEventV1Discard), id).Get(ctx, nil)
			if err != nil {
				log.Fatalln("Failure executing EventDiscard", err)
			}
			discarded = true
		})

		updateChangesCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1UpdateChanges))
		discardSubmitOrAddChangesSelector.AddReceive(updateChangesCh, func(ch workflow.ReceiveChannel, more bool) {
			changes := make([]model.EntityChange, 0)
			ch.Receive(ctx, &changes)
			var result interface{}
			err := workflow.ExecuteActivity(ctx, v1.ChangeRequestWorkflowEventV1UpdateChanges, changes).Get(ctx, &result)
			if err != nil {
				log.Fatalln("Failure executing EventUpdateChanges", err)
			}
			// TODO: Check if the changes are updated in the activities object
		})

		submitForReviewCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1SubmitForReview))
		discardSubmitOrAddChangesSelector.AddReceive(submitForReviewCh, func(ch workflow.ReceiveChannel, more bool) {
			ch.Receive(ctx, nil)
			workflow.GetLogger(ctx).Info("Submitting change request for review", "request", id)
			err := workflow.ExecuteActivity(ctx, string(v1.ChangeRequestWorkflowEventV1SubmitForReview), id).Get(ctx, nil)
			if err != nil {
				log.Fatalln("Failure executing EventSubmitForReview", err)
			}
			submitted = true
		})

		discardSubmitOrAddChangesSelector.Select(ctx)
		if discarded == true {
			return nil
		}
		if submitted == true {
			break
		}
	}

	reviewReceived := false

	for {
		signalSelector := workflow.NewSelector(ctx)

		// 3. Intermediate Approval Callbacks
		intermediateApprovalCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1IntermediateReviewCallbacks))
		signalSelector.AddReceive(intermediateApprovalCh, func(ch workflow.ReceiveChannel, more bool) {
			workflow.GetLogger(ctx).Info("Received intermediate approval signal")
			approver := ""
			ch.Receive(ctx, &approver)
			workflow.GetLogger(ctx).Info("Approver", "approver", approver)

			workflow.GetLogger(ctx).Info("Executing intermediate approval callback")
			var result interface{}
			err := workflow.ExecuteActivity(ctx, string(v1.ChangeRequestWorkflowEventV1IntermediateReviewCallbacks), id, approver).Get(ctx, &result)
			if err != nil {
				log.Fatalln("Failure executing EventIntermediateReviewCallback", err)
			}
		})

		// 4. Final Review - Approval/Rejection
		finalReviewCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1FinalReview))
		signalSelector.AddReceive(finalReviewCh, func(ch workflow.ReceiveChannel, more bool) {

			workflow.GetLogger(ctx).Info("Received final review signal")
			verdict := ""
			ch.Receive(ctx, &verdict)
			workflow.GetLogger(ctx).Info("Verdict", "verdict", verdict)

			workflow.GetLogger(ctx).Info("Executing final review callback")
			var result interface{}
			err := workflow.ExecuteActivity(ctx, string(v1.ChangeRequestWorkflowEventV1FinalReview), id, verdict).Get(ctx, &result)
			if err != nil {
				log.Fatalln("Failure executing EventFinalReviewCallback", err)
			}
			reviewReceived = true
		})
		signalSelector.Select(ctx)
		if reviewReceived == true {
			break
		}
	}
	return nil
	// @@@SNIPEND
}
