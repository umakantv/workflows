# workflows
This repository contains the tutorial code for implementing workflows in Go using the Cadence and Temporal Workflow libraries.

## Temporal

Docs: https://docs.temporal.io/docs/go/workflows

## Cadence

Docs: https://cadenceworkflow.io/docs/get-started/

## How to run

1. Update the `DB_PATH`  value in `common/repo/change_request_repo.go` with root directory of the project.
2. Run `go run worker/main.go` to start the worker.
3. Run `go run script/main.go` to start the trigger events.

Modify the `script/main.go` to trigger different events.