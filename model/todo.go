package model

//defines the data model for the todo application
type TodoItem struct {
	ID          int
	Description string
	Status      string
}
type APITodoRequest struct {
	Description string
	Status      string
}
