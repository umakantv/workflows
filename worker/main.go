package main

import (
	"github.com/umakantv/workflows/temporal"
)

func main() {
	// Run the worker to execute the workflows and activities
	temporal.InitWorker()
}
