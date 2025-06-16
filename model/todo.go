package model

type TodoItem struct {
	ID          int
	Description string
	Status      string
}
type APITodoRequest struct {
	Description string
	Status      string
}
