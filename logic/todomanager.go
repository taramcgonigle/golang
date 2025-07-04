package logic

import (
	"practice/assignments/model"
	"practice/assignments/storage"
	"sync"
	"testing"
)

// commandType defines the type of operation sent to the TodoManager
type commandType int

const (
	getTodos commandType = iota
	addTodo
	updateTodo
	deleteTodo
)

// command represents a message sent to the TodoManager actor
type command struct {
	typ      commandType
	item     model.TodoItem // used by addTodo
	payload  model.TodoItem // used for updateTodo (can be reused for delete later)
	index    int            // used for deleteTodo
	response chan any       // channel to send back the response
}

// TodoManager is the actor that manages to-do list access safely via channels
type TodoManager struct {
	cmdChan chan command
}

// NewTodoManager starts the actor goroutine and initializes with stored todos
func NewTodoManager() *TodoManager {
	manager := &TodoManager{cmdChan: make(chan command)}
	todos := storage.LoadTodos() // Load todos at startup

	//actor loop that handles the commands
	go func() {
		for cmd := range manager.cmdChan {
			switch cmd.typ {

			case getTodos:
				// Return a copy of the todos list
				cmd.response <- append([]model.TodoItem{}, todos...)

			case addTodo:
				todos = append(todos, cmd.item)
				storage.SaveTodos(todos)
				cmd.response <- struct{}{} // acknowledge add

			case updateTodo:
				todo := cmd.payload
				updated := false
				for i, item := range todos {
					if item.ID == todo.ID {
						// Apply updates by ID match
						todos[i].Description = todo.Description
						todos[i].Status = todo.Status
						storage.SaveTodos(todos)
						updated = true
						break
					}
				}
				cmd.response <- updated // true if update successful

			case deleteTodo:
				id := cmd.index // This is the ID of the item
				found := false
				for i, t := range todos {
					if t.ID == id {
						todos = append(todos[:i], todos[i+1:]...) // Remove item by ID
						storage.SaveTodos(todos)
						found = true
						break
					}
				}
				cmd.response <- found // true if delete succeeded

			}
		}
	}()

	return manager
}

// GetAll retrieves a copy of the to-do list
func (m *TodoManager) GetAll() []model.TodoItem {
	resp := make(chan any)
	m.cmdChan <- command{typ: getTodos, response: resp}
	return (<-resp).([]model.TodoItem)
}

// Add adds a new item to the list
func (m *TodoManager) Add(item model.TodoItem) {
	resp := make(chan any)
	m.cmdChan <- command{typ: addTodo, item: item, response: resp}
	<-resp
}

// Update modifies an existing to-do item
func (m *TodoManager) Update(todo model.TodoItem) {
	resp := make(chan any)
	m.cmdChan <- command{typ: updateTodo, payload: todo, response: resp}
	<-resp
}

// Deletes an existing to-do item
func (m *TodoManager) Delete(index int) {
	resp := make(chan any)
	m.cmdChan <- command{typ: deleteTodo, index: index, response: resp}
	<-resp
}

// ensures that the manager handles concurrent access (add + update) safely using goroutines
func TestTodoManager_ConcurrentAddAndUpdate(t *testing.T) {
	t.Parallel()

	const goroutines = 20
	var wg sync.WaitGroup

	// Use a fresh manager for this test to avoid interference
	m := NewTodoManager()

	// Concurrently add todos with unique IDs
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Add(model.TodoItem{
				ID:          10000 + i, // ensure unique IDs
				Description: "desc",
				Status:      "not started",
			})
		}(i)
	}
	wg.Wait()

	// Concurrently update todos
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Update(model.TodoItem{
				ID:          10000 + i,
				Description: "updated",
				Status:      "completed",
			})
		}(i)
	}
	wg.Wait()

	// Validate all todos items were updated
	todos := m.GetAll()
	if len(todos) != goroutines {
		t.Fatalf("expected %d todos, got %d", goroutines, len(todos))
	}
	for _, todo := range todos {
		if todo.Description != "updated" || todo.Status != "completed" {
			t.Errorf("todo not updated correctly: %+v", todo)
		}
	}
}

// NewInMemoryTodoManager creates a manager that does not persist to disk (for tests)
func NewInMemoryTodoManager() *TodoManager {
	manager := &TodoManager{cmdChan: make(chan command)}
	var todos []model.TodoItem

	go func() {
		for cmd := range manager.cmdChan {
			switch cmd.typ {
			case getTodos:
				cmd.response <- append([]model.TodoItem(nil), todos...)
			case addTodo:
				todos = append(todos, cmd.item)
				cmd.response <- struct{}{}
			case updateTodo:
				todo := cmd.payload
				updated := false
				for i, item := range todos {
					if item.ID == todo.ID {
						todos[i].Description = todo.Description
						todos[i].Status = todo.Status
						updated = true
						break
					}
				}
				cmd.response <- updated
			case deleteTodo:
				id := cmd.index
				found := false
				for i, t := range todos {
					if t.ID == id {
						todos = append(todos[:i], todos[i+1:]...)
						found = true
						break
					}
				}
				cmd.response <- found
			}
		}
	}()
	return manager
}
