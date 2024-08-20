package activities

import (
	"context"
	"errors"
	"log"

	"github.com/umakantv/workflows/common/model"
	"github.com/umakantv/workflows/common/repo"
)

func (c *ChangeRequestActivitiesV1) EventFinalReviewCallback(ctx context.Context, id, verdict string) error {
	// @@@SNIPBEGIN temporal/change_request/v1/activities/UpdateChanges
	// This activity submits the change request for review.

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Printf("Error: Failed to get change request %s\n", id)
		return err
	}
	log.Printf("Final Review Callback for %s - %s\n", changeRequest.ID, verdict)

	if changeRequest.Status != model.ChangeRequestStatusSubmitted {
		return errors.New("change request is not in submitted state")
	}

	// Merge the changes to the entity
	log.Println("Applying change")

	// Send email to the creator
	if verdict == "approved" {
		log.Printf("Change request has been approved")
	} else {
		log.Printf("Change request has been rejected")
	}

	if verdict == "approved" {
		changeRequest.Status = model.ChangeRequestStatusApproved
	} else {
		changeRequest.Status = model.ChangeRequestStatusRejected
	}
	err = repo.GetChangeRequestRepo().UpdateChangeRequest(changeRequest)
	if err != nil {
		log.Printf("Error: Failed to update change request %s\n", id)
		return err
	}

	return nil
	// @@@SNIPEND
}
