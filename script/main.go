package main

import (
	"log"
	"time"

	"github.com/umakantv/workflows/common/app"
	"github.com/umakantv/workflows/temporal"
)

func main() {

	temporal.Init()

	// Test the workflow with temporal
	id := app.GetCore().InitializeChangeRequest()

	time.Sleep(5 * time.Second)
	app.GetCore().SubmitForReview(id)

	time.Sleep(5 * time.Second)
	log.Println("Change request status", app.GetCore().GetChangeRequestStatus(id))

	time.Sleep(5 * time.Second)
	app.GetCore().ApproveChangeRequest(id)

	time.Sleep(3 * time.Second)
	log.Println("Change request status", app.GetCore().GetChangeRequestStatus(id))

}
