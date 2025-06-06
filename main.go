package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	//"io/ioutil"
)

// defining strct to represent a single to-do list
type TodoItem struct {
	ID          int    `json:"id"`          //unique ID
	Description string `json:"description"` //task description
	Status      string `json:"status"`      // Statuses: "not started", "started", "completed"
}

// Storing of to-do list
var todoFile = "todo.json"

func loadTodos() []TodoItem {
	var todos []TodoItem

	//tries to read the file
	file, err := ioutil.ReadFile(todoFile)
	if err != nil {
		//if no file exists then return empty list
		return todos
	}

	//convert JSON to Go slice of struts
	_ = json.Unmarshal(file, &todos)
	return todos
}

func saveTodos(todos []TodoItem) {
	//Convert Go slice of struts to JSON
	data, _ := json.MarshalIndent(todos, "", "  ")
	//write JSON to file
	_ = ioutil.WriteFile(todoFile, data, 0644)
}

func main() {
	//CLI calls
	add := flag.String("add", "", "Add a new todo item")
	update := flag.Int("update", -1, "Update the description of a todo item by ID")
	newDesc := flag.String("desc", "", "New description for update")
	deleteID := flag.Int("delete", -1, "Delete a todo item by ID")
	status := flag.String("status", "", "Set status (not started, started, completed)")
	list := flag.Bool("list", false, "List all to-do items")

	//parse flags from command line input
	flag.Parse()

	//load existing to-do list
	todos := loadTodos()

	//add item
	if *add != "" {
		newID := len(todos) + 1
		todos = append(todos, TodoItem{
			ID:          newID,
			Description: *add,
			Status:      "not started",
		})
		fmt.Println("Added new todo item.")
	}

	//update an item of the list
	if *update != -1 && *newDesc != "" {
		for i := range todos {
			if todos[i].ID == *update {
				todos[i].Description = *newDesc
				fmt.Println("Updated description.")
			}
		}
	}

	//update status of an item
	if *status != "" && *update != -1 {
		for i := range todos {
			if todos[i].ID == *update {
				todos[i].Status = *status
				fmt.Println("Updated status.")
			}
		}
	}

	//delete an item
	if *deleteID != -1 {
		for i := range todos {
			if todos[i].ID == *deleteID {
				todos = append(todos[:i], todos[i+1:]...)
				fmt.Println("Deleted item.")
				break
			}
		}
	}

	//list all iems
	if *list || *add != "" || *update != -1 || *deleteID != -1 {
		fmt.Println("Current To-Do List:")
		for _, item := range todos {
			fmt.Printf("%d: %s [%s]\n", item.ID, item.Description, item.Status)
		}
	}

	//save the list back to the file
	saveTodos(todos)
}
