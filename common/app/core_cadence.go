package app

import (
	"context"
	"log"
	"time"

	"github.com/umakantv/workflows/common/repo"
	"github.com/umakantv/workflows/temporal"
	v1 "github.com/umakantv/workflows/temporal/change_request/v1"
	"github.com/umakantv/workflows/temporal/change_request/v1/workflow"
	"go.temporal.io/sdk/client" // TODO: Switch to cadence client
)

type CoreCadence struct {
	cadenceClient client.Client
}

func NewCoreCadence() *CoreCadence {
	return &CoreCadence{
		cadenceClient: temporal.GetClient(), // TODO: Switch to cadence client
	}
}

func (core *CoreCadence) InitializeChangeRequest() string {

	id, err := repo.GetChangeRequestRepo().InitiateChangeRequest()
	if err != nil {
		log.Fatalln("Unable to initiate change request", err)
	}

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Fatalln("Unable to get change request", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:                 "Change_request_v1_" + id,
		TaskQueue:          v1.ChangeRequestTaskQueueGoV1,
		WorkflowRunTimeout: time.Hour * 24 * 10,
	}

	we, err := core.cadenceClient.ExecuteWorkflow(
		context.Background(), workflowOptions, workflow.ChangeRequestWorkflowV1, id)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	changeRequest.WorkflowID = we.GetID()
	changeRequest.RunID = we.GetRunID()

	err = repo.GetChangeRequestRepo().UpdateChangeRequest(changeRequest)
	if err != nil {
		log.Fatalln("Unable to update change request", err)
	}

	return id
}

func (core *CoreCadence) GetChangeRequestStatus(id string) string {

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Fatalln("Unable to get change request", err)
	}

	return string(changeRequest.Status)
}

func (core *CoreCadence) SubmitForReview(id string) {

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Fatalln("Unable to get change request", err)
	}

	err = core.cadenceClient.SignalWorkflow(
		context.Background(), changeRequest.WorkflowID, changeRequest.RunID,
		string(v1.ChangeRequestWorkflowEventV1SubmitForReview), nil)
	if err != nil {
		log.Fatalln("Unable to send workflow signal", err)
	}
	log.Println("Sent signal to workflow for submitting changes for review", "WorkflowID", changeRequest.WorkflowID)
}

func (core *CoreCadence) ApproveChangeRequest(id string) {
	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Fatalln("Unable to get change request", err)
	}

	err = core.cadenceClient.SignalWorkflow(
		context.Background(), changeRequest.WorkflowID, changeRequest.RunID,
		string(v1.ChangeRequestWorkflowEventV1FinalReview), "approved")
	if err != nil {
		log.Fatalln("Unable to send workflow signal", err)
	}
	log.Println("Sent signal to workflow for approving changes", "WorkflowID", changeRequest.WorkflowID)
}

func (core *CoreCadence) RejectChangeRequest(id string) {

	changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	if err != nil {
		log.Fatalln("Unable to get change request", err)
	}

	err = core.cadenceClient.SignalWorkflow(
		context.Background(), changeRequest.WorkflowID, changeRequest.RunID,
		string(v1.ChangeRequestWorkflowEventV1FinalReview), "rejected")
	if err != nil {
		log.Fatalln("Unable to send workflow signal", err)
	}
	log.Println("Sent signal to workflow for rejecting changes", "WorkflowID", changeRequest.WorkflowID)
}
