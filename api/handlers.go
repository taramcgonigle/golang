package api

// defines the API handlers for the todo application, and interacts with the concurrency todomanager

import (
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"practice/assignments/logic"
	"practice/assignments/model"
	"practice/assignments/trace"
	"strconv"
)

// global instance for all handlers to access todo list via concurrency
var manager *logic.TodoManager

func init() { //initialize the todo manager at package load time
	manager = logic.NewTodoManager()
}

// Place getTemplatePath here
func getTemplatePath(name string) string {
	// Try local templates/ first
	path := filepath.Join("templates", name)
	if _, err := os.Stat(path); err == nil {
		return path
	}
	// Try parent directory (for tests run from api/)
	path = filepath.Join("..", "templates", name)
	if _, err := os.Stat(path); err == nil {
		return path
	}
	// Fallback: original
	return filepath.Join("templates", name)
}

// get handler serves the homepage which displays the list of todos
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context()) //log the request trace ID
	slog.Info("Serving homepage", "traceID", traceID)

	//todos := storage.LoadTodos() //load all current todos from json file

	todos := manager.GetAll() //use actor managed state

	//render the homepage template + passing the todos
	tmpl := template.Must(template.ParseFiles(getTemplatePath("home.html")))
	err := tmpl.Execute(w, todos)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
	}
}

// CreateFormHandler renders the form used to create a new todo item
func CreateFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(getTemplatePath("create.html")))
	tmpl.Execute(w, nil)
}

// CreateHandler handles POST /create. It parses form data, validates inputs, creates a new todo item, and redirects to the home page.
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context())

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	desc := r.FormValue("description")
	status := r.FormValue("status")

	if desc == "" || status == "" { // make sure both fields are provided
		http.Error(w, "Missing description or status", http.StatusBadRequest)
		return
	}

	// Get current todos to calculate next ID
	todos := manager.GetAll()
	newID := 0
	if len(todos) > 0 {
		newID = todos[len(todos)-1].ID + 1
	}

	// Construct and add the new item
	todo := model.TodoItem{
		ID:          newID,
		Description: desc,
		Status:      status,
	}

	// Send to actor
	manager.Add(todo)

	slog.Info("Todo created", "traceID", traceID, "ID", todo.ID, "description", todo.Description)

	// Redirect back to home after successful creation
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Handles /get and returns the list of todo items as a JSON response
func GetHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context()) //gets trace ID for logging
	slog.Info("Returning JSON todo list", "traceID", traceID)

	//todos := storage.LoadTodos()

	// Use the concurrency-safe manager to get the current todo list
	todos := manager.GetAll()

	w.Header().Set("Content-Type", "application/json") //return the todos as json
	json.NewEncoder(w).Encode(todos)
}

// Renders the edit form pre-filled with data for a specific todo item, identified by ID.
func EditFormHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context())
	slog.Info("Edit form request", "traceID", traceID)

	idStr := r.FormValue("id") // works for both GET and POST
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	todos := manager.GetAll()

	// Search for the todo with the given ID
	var todo model.TodoItem
	found := false
	for _, t := range todos {
		if t.ID == id {
			todo = t
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "To-do item not found", http.StatusNotFound)
		return
	}

	// Render the form with the existing item preloaded
	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	if err := tmpl.Execute(w, todo); err != nil {
		http.Error(w, "Error rendering edit page", http.StatusInternalServerError)
	}
}

// handles POST /update to modify an existing todo item.
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	// Extract values from the form
	idStr := r.FormValue("id")
	desc := r.FormValue("description")
	status := r.FormValue("status")

	traceID := trace.GetTraceID(r.Context())
	slog.Info("Update requested", "traceID", traceID, "id", idStr, "description", desc, "status", status)

	// Convert ID from string to int
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Build the updated TodoItem to submit to the manager
	todo := model.TodoItem{
		ID:          id,
		Description: desc,
		Status:      status,
	}

	// Use the manager to update the item (concurrency-safe)
	manager.Update(todo)

	slog.Info("Todo updated via manager", "traceID", traceID, "id", id)

	// Redirect back to homepage
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// handles GET /delete?id=ID to remove a todo item. Reads ID from query parameters.
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	traceID := trace.GetTraceID(r.Context())

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing 'id' parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	manager.Delete(id)

	slog.Info("Deleted todo", "traceID", traceID, "id", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
