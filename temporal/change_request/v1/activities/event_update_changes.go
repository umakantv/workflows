package activities

import (
	"context"
	"errors"
	"log"

	"github.com/umakantv/workflows/common/model"
	"github.com/umakantv/workflows/common/repo"
)

func (c *ChangeRequestActivitiesV1) EventUpdateChanges(ctx context.Context,
	id string, change model.EntityChange) error {
	// @@@SNIPBEGIN temporal/change_request/v1/activities/UpdateChanges
	// This activity updates the changes made to an entity.
	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Printf("Error: Failed to get change request %s\n", id)
		return err
	}
	log.Println("Updating changes for change request", changeRequest.ID)

	if changeRequest.Status != model.ChangeRequestStatusDraft {
		return errors.New("change request is not in draft state")
	}

	// changeRequest.Change = change

	err = repo.GetChangeRequestRepo().UpdateChangeRequest(changeRequest)
	if err != nil {
		log.Printf("Error: Failed to update change request %s\n", id)
		return err
	}

	return nil
	// @@@SNIPEND
}
