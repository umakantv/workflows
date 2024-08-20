package workflow

import (
	"log"
	"time"

	"github.com/umakantv/workflows/common/model"
	"github.com/umakantv/workflows/common/repo"
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

	// Register query handlers
	err := workflow.SetQueryHandler(ctx, string(v1.ChangeRequestWorkflowQueryV1GetStatus),
		func(input []byte) (string, error) {
			changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
			if err != nil {
				log.Fatalln("Failure getting change request", err)
				return "", err
			}
			return string(changeRequest.Status), nil
		})
	if err != nil {
		log.Fatalln("Failure setting query handler", err)
		return err
	}

	// Signals we need to handle
	// 1. Update Changes
	// 2. Discard
	// 3. Submit for Review
	// 4. Intermediate Approval Callbacks
	// 5. Final Review - Approval/Rejection

	// Define signal channels
	// 1. Update Changes
	updateChangesSelector := workflow.NewSelector(ctx)
	updateChangesCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1UpdateChanges))
	updateChangesSelector.AddReceive(updateChangesCh, func(ch workflow.ReceiveChannel, more bool) {
		changes := make([]model.EntityChange, 0)
		ch.Receive(ctx, &changes)
		var result interface{}
		err := workflow.ExecuteActivity(ctx, v1.ChangeRequestWorkflowEventV1UpdateChanges, changes).Get(ctx, &result)
		if err != nil {
			log.Fatalln("Failure executing EventUpdateChanges", err)
		}
		// TODO: Check if the changes are updated in the activities object
	})

	// 2. Discard
	discardSelector := workflow.NewSelector(ctx)
	discardCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1Discard))
	discardSelector.AddReceive(discardCh, func(ch workflow.ReceiveChannel, more bool) {
		var result interface{}
		workflow.GetLogger(ctx).Info("Discarding change request", "request", id)
		err := workflow.ExecuteActivity(ctx, string(v1.ChangeRequestWorkflowEventV1Discard), id).Get(ctx, &result)
		if err != nil {
			log.Fatalln("Failure executing EventDiscard", err)
		}
	})
	// discardSelector.Select(ctx)

	// 3. Submit for Review
	submitForReviewSelector := workflow.NewSelector(ctx)
	submitForReviewCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1SubmitForReview))
	submitForReviewSelector.AddReceive(submitForReviewCh, func(ch workflow.ReceiveChannel, more bool) {
		workflow.GetLogger(ctx).Info("Submitting change request for review", "request", id)
		err := workflow.ExecuteActivity(ctx, string(v1.ChangeRequestWorkflowEventV1SubmitForReview), id).Get(ctx, nil)
		if err != nil {
			log.Fatalln("Failure executing EventSubmitForReview", err)
		}
	})
	submitForReviewSelector.Select(ctx)

	// 3. Intermediate Approval Callbacks
	intermediateApprovalSelector := workflow.NewSelector(ctx)
	intermediateApprovalCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1IntermediateReviewCallbacks))
	intermediateApprovalSelector.AddReceive(intermediateApprovalCh, func(ch workflow.ReceiveChannel, more bool) {

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
	finalReviewSelector := workflow.NewSelector(ctx)
	finalReviewCh := workflow.GetSignalChannel(ctx, string(v1.ChangeRequestWorkflowEventV1FinalReview))
	finalReviewSelector.AddReceive(finalReviewCh, func(ch workflow.ReceiveChannel, more bool) {

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
	})
	finalReviewSelector.Select(ctx)

	// Start the workflow
	// Busy waiting for the change request status to reach terminal state
	workflow.AwaitWithTimeout(ctx, 365*24*time.Hour, func() bool {
		changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
		if err != nil {
			workflow.GetLogger(ctx).Error("Failed to get change request", "error", err)
			return false
		}
		workflow.GetLogger(ctx).Info("Checking change request status", "status", changeRequest.Status)
		if changeRequest.Status == model.ChangeRequestStatusApproved ||
			changeRequest.Status == model.ChangeRequestStatusRejected ||
			changeRequest.Status == model.ChangeRequestStatusDiscarded {
			return true
		}
		return false
	})

	return nil
	// @@@SNIPEND
}
