package jsonplaceholder

import (
	_ "embed"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"slices"
	"strconv"
	"testing"
)

//go:embed testdata/tododb.json
var tododbJSON []byte

var (
	httpServer *httptest.Server
	tododb     []Todo

	testdata = map[int64]Todo{
		1: {ID: 1, Title: "Todo 1", Completed: false, UserID: 123},
		2: {ID: 2, Title: "Todo 2", Completed: false, UserID: 123},
	}
)

func TestMain(m *testing.M) {
	// setup
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
	err := json.Unmarshal(tododbJSON, &tododb)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /todos/{id}", func(w http.ResponseWriter, r *http.Request) {
		idstr := r.PathValue("id")
		id, err := strconv.ParseInt(idstr, 10, 64)
		if err != nil {
			log.Error("get todos", "error", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		idx := slices.IndexFunc(tododb, func(t Todo) bool {
			return t.ID == id
		})

		if idx < 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(tododb[idx]); err != nil {
			log.Error("cannot write response body", "error", err)
		}
	})

	httpServer = httptest.NewServer(mux)
	exit := m.Run()
	// teardown
	httpServer.Close()
	os.Exit(exit)
}

// DoerFunc is adapter that convert ordinary function into Doer.
type DoerFunc func(req *http.Request) (*http.Response, error)

func (f DoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestClient_GetTodo(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		c := NewClient(DoerFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != "GET" {
				return nil, errors.New("method not allowed")
			}
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusNotFound)
			return rec.Result(), nil
		}), "")

		ctx := t.Context()
		_, err := c.GetTodo(ctx, 1)
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	})
	t.Run("found", func(t *testing.T) {
		expectedBody := `
			{
				"userId": 1,
				"id": 1,
				"title": "delectus aut autem",
				"completed": false
			  }
			`
		c := NewClient(DoerFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != "GET" {
				return nil, errors.New("method not allowed")
			}
			rec := httptest.NewRecorder()
			rec.WriteHeader(http.StatusOK)
			_, err := rec.WriteString(expectedBody)
			if err != nil {
				t.Fatalf("cannot write response body: %s", err.Error())
			}
			return rec.Result(), nil
		}), "")

		ctx := t.Context()
		todo, err := c.GetTodo(ctx, 1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if todo.ID != 1 {
			t.Errorf("unexpected todo id: actual %d", todo.ID)
		}

		if todo.Title != "delectus aut autem" {
			t.Errorf("unexpected title: actual %q", todo.Title)
		}
	})
}

func TestClientWithTestServer_GetTodo(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		c := NewClient(httpServer.Client(), httpServer.URL)
		ctx := t.Context()
		_, err := c.GetTodo(ctx, -1)
		if !errors.Is(err, ErrTodoNotFound) {
			t.Fatalf("unexpected error: %s", err.Error())
		}
	})
	t.Run("found", func(t *testing.T) {
		c := NewClient(httpServer.Client(), httpServer.URL)

		ctx := t.Context()
		todo, err := c.GetTodo(ctx, 1)
		if err != nil {
			t.Fatalf("unexpected error: %s", err.Error())
		}

		if todo.ID != 1 {
			t.Errorf("unexpected todo id: actual %d", todo.ID)
		}

		idx := slices.IndexFunc(tododb, func(t Todo) bool { return t.ID == 1 })
		if todo.Title != tododb[idx].Title {
			t.Errorf("unexpected title: actual %q", todo.Title)
		}
	})
}
