# golang

1. A simple command line to-do list application. It allows you to create, view, update, and delete tasks with persistent storage using a local JSON file.

Features:
- Add new to-do items via command line flags --> go run main.go -add "Buy groceries"

- View the complete list of to-dos with their statuses --> go run main.go -list

- Update task descriptions --> go run main.go -update [id] -desc "to do item"
        and statuses       -->go run main.go -update [id] -status [completed, started, "not started"]

- Delete items by ID --> go run main.go -delete [id]

- Automatically saves to and loads from disk {todo.json}. To reset the app --> rm todo.json

Running Tests:
Run all tests with verbose output -->  go test ./... -v

To view test coverage --> go test ./... -coverprofile=coverage.out