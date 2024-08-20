package temporal

import (
	"log"

	v1 "github.com/umakantv/workflows/temporal/change_request/v1"
	"github.com/umakantv/workflows/temporal/change_request/v1/activities"
	"github.com/umakantv/workflows/temporal/change_request/v1/workflow"
	"go.temporal.io/sdk/worker"
)

func InitWorker() {

	InitClient()
	// The client and Worker are heavyweight objects that should be created once per process.
	c := GetClient()

	w := worker.New(c, v1.ChangeRequestTaskQueueGoV1, worker.Options{})

	w.RegisterWorkflow(workflow.ChangeRequestWorkflowV1)
	w.RegisterActivity(&activities.ChangeRequestActivitiesV1{})

	err := w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
