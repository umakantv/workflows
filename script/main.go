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
	id := app.InitializeChangeRequest()

	// id, err := repo.GetChangeRequestRepo().InitiateChangeRequest()
	// if err != nil {
	// 	log.Fatalln("Unable to initiate change request", err)
	// }
	// changeRequest, err := repo.GetChangeRequestRepo().GetChangeRequest(id)
	// if err != nil {
	// 	log.Fatalln("Unable to get change request", err)
	// }
	// log.Println("Change request", changeRequest)

	time.Sleep(5 * time.Second)
	app.SubmitForReview(id)

	time.Sleep(5 * time.Second)
	log.Println("Change request status", app.GetChangeRequestStatus(id))

	time.Sleep(5 * time.Second)
	app.ApproveChangeRequest(id)

	time.Sleep(3 * time.Second)
	log.Println("Change request status", app.GetChangeRequestStatus(id))

}
