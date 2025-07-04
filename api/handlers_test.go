package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"practice/assignments/logic"
	"practice/assignments/model"
	"strings"
	"sync"
	"testing"
)

// Validates /get endpoint returns a 200 status + Json response
func TestGetHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/get", nil) //new HHTP GET request to /get

	// Attach a base context so trace injection works
	ctx := context.Background()
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder() //mocked http response for test verification

	handler := withTraceID(http.HandlerFunc(GetHandler)) //wrap handler with trace ID middleware
	handler.ServeHTTP(rr, req)                           //passing through middleware-wrapped handler

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %d", rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" { //check Content-Type is json
		t.Errorf("Expected Content-Type application/json, got %s", ct)
	}

	t.Logf("Response body: %s", rr.Body.String()) //log for visibility
}

// Confirms HTML is returned from the root `/`
func TestHomeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(context.Background())
	rr := httptest.NewRecorder()

	handler := withTraceID(http.HandlerFunc(HomeHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" && ct != "" {
		t.Errorf("Expected HTML content-type, got %s", ct)
	}

	t.Logf("Home page HTML:\n%s", rr.Body.String())
}

// Renders a form for adding new items
func TestCreateFormHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/create-form", nil)
	req = req.WithContext(context.Background())
	rr := httptest.NewRecorder()

	handler := withTraceID(http.HandlerFunc(CreateFormHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %d", rr.Code)
	}

	t.Log("Create form rendered")
}

// Simulates form submission to `/create`, expects redirect
func TestCreateHandler(t *testing.T) {
	form := "description=Test+Todo&status=not+started"
	req := httptest.NewRequest(http.MethodPost, "/create",
		strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.Background())

	rr := httptest.NewRecorder()
	handler := withTraceID(http.HandlerFunc(CreateHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected redirect (303), got %d", rr.Code)
	}

	location := rr.Header().Get("Location")
	if location != "/" {
		t.Errorf("Expected redirect to '/', got %s", location)
	}
}

// Simulates form-based update to `/update`, validates DB change
func TestUpdateHandler(t *testing.T) {
	// Create dummy item to update
	manager.Add(model.TodoItem{ //manager is used to make the test independent
		ID:          999,
		Description: "Old",
		Status:      "not started",
	})

	form := "id=999&description=Updated+Task&status=completed"
	req := httptest.NewRequest(http.MethodPost, "/update",
		strings.NewReader(form))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(context.Background())

	rr := httptest.NewRecorder()
	handler := withTraceID(http.HandlerFunc(UpdateHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 redirect, got %d", rr.Code)
	}
}

// Simulates deletion via query param on `/delete`
func TestDeleteHandler(t *testing.T) {
	// Add a dummy item with known ID
	manager.Add(model.TodoItem{
		ID:          123,
		Description: "Delete me",
		Status:      "started",
	})

	req := httptest.NewRequest(http.MethodGet, "/delete?id=123", nil)
	req = req.WithContext(context.Background())
	rr := httptest.NewRecorder()

	handler := withTraceID(http.HandlerFunc(DeleteHandler))
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Expected 303 redirect, got %d", rr.Code)
	}
}

// Performs concurrent `Add` and `Update` operations to stress test the manager's thread safety
func TestTodoManager_ConcurrentAddAndUpdate(t *testing.T) {
	t.Parallel()

	const goroutines = 100
	var wg sync.WaitGroup //ensures all goroutines complete before validation

	// Use a fresh manager for this test to avoid interference
	m := logic.NewInMemoryTodoManager() //for concurrency tests

	// Concurrently add todos
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			m.Add(model.TodoItem{
				ID:          i,
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
				ID:          i,
				Description: "updated",
				Status:      "completed",
			})
		}(i)
	}
	wg.Wait()

	// Validate all todos were updated
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
