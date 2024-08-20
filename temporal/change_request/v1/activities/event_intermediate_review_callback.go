package activities

import (
	"context"
	"errors"
	"log"

	"github.com/umakantv/workflows/common/model"
	"github.com/umakantv/workflows/common/repo"
)

func (c *ChangeRequestActivitiesV1) EventIntermediateReviewCallback(ctx context.Context, id, approver string) error {
	// @@@SNIPBEGIN temporal/change_request/v1/activities/UpdateChanges
	// This activity submits the change request for review.

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Printf("Error: Failed to get change request %s\n", id)
		return err
	}
	log.Printf("Change Request Review Callback for %s\n", changeRequest.ID)

	if changeRequest.Status != model.ChangeRequestStatusSubmitted {
		return errors.New("change request is not in submitted state")
	}

	if changeRequest.Approvers == "" {
		changeRequest.Approvers = approver
	} else {
		changeRequest.Approvers += ", " + approver
	}
	repo.GetChangeRequestRepo().UpdateChangeRequest(changeRequest)

	// Send email to the creator
	log.Printf("%s has approved the change request", approver)

	return nil
	// @@@SNIPEND
}
