package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/umakantv/workflows/temporal"
)

func Start() {

	temporal.Init()

	muxRouter := mux.NewRouter()

	// register routes
	muxRouter.HandleFunc("/initiate", initiateChangeRequestHandler).Methods("POST")
	muxRouter.HandleFunc("/submit", submitForReviewHandler).Methods("POST")
	muxRouter.HandleFunc("/status", getChangeRequestStatusHandler).Methods("GET")
	muxRouter.HandleFunc("/approve", approveChangeRequestHandler).Methods("POST")
	muxRouter.HandleFunc("/reject", rejectChangeRequestHandler).Methods("POST")

	// start server
	err := http.ListenAndServe(":3010", muxRouter)
	if err != nil {
		log.Println("Error starting server", err)
		return
	}

	log.Println("Server started at port 3010")
}
