package activities

import (
	"context"
	"errors"
	"log"

	"github.com/umakantv/workflows/common/model"
	"github.com/umakantv/workflows/common/repo"
)

func (c *ChangeRequestActivitiesV1) EventDiscard(ctx context.Context, id string) error {
	// @@@SNIPBEGIN temporal/change_request/v1/activities/UpdateChanges
	// This activity submits the change request for review.

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Printf("Error: Failed to get change request %s\n", id)
		return err
	}
	log.Printf("Discarding change request %s\n", changeRequest.ID)

	if changeRequest.Status != model.ChangeRequestStatusDraft {
		return errors.New("change request is not in draft state")
	}

	changeRequest.Status = model.ChangeRequestStatusDiscarded
	err = repo.GetChangeRequestRepo().UpdateChangeRequest(changeRequest)
	if err != nil {
		log.Printf("Error: Failed to update change request %s\n", id)
		return err
	}

	return nil
	// @@@SNIPEND
}
