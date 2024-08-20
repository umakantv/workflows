package app

type Core interface {
	InitializeChangeRequest() string
	GetChangeRequestStatus(id string) string
	SubmitForReview(id string)
	ApproveChangeRequest(id string)
	RejectChangeRequest(id string)
}

var coreInstance Core

func GetCore() Core {
	return coreInstance
}

func SetCore(core Core) {
	coreInstance = core
}
