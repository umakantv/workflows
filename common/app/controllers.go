package app

import (
	"net/http"

	"github.com/umakantv/workflows/common/repo"
)

func initiateChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	InitializeChangeRequest()

	w.Write([]byte("Change request initiated"))
	w.WriteHeader(http.StatusOK)
}

func getChangeRequestStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	request, err := repo.GetChangeRequestRepo().GetChangeRequest("1")
	if err != nil {
		http.Error(w, "Failed to get change request", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(request.Status))
	w.WriteHeader(http.StatusOK)
}

func submitForReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	SubmitForReview("1")

	w.Write([]byte("Change request submitted for review"))
}

func approveChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	ApproveChangeRequest("1")
}

func rejectChangeRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	RejectChangeRequest("1")
}
