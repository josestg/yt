package jsonplaceholder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrTodoNotFound = errors.New("jsonplaceholder: todo not found")
)

type Todo struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"userId"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	http         Doer
	listTodoPath string
}

func NewClient(c Doer, baseURL string) *Client {
	return &Client{
		http:         c,
		listTodoPath: strings.Join([]string{baseURL, "todos"}, "/"),
	}
}

func (c *Client) ListTodo(ctx context.Context) ([]Todo, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.listTodoPath, nil)
	if err != nil {
		return nil, fmt.Errorf("jsonplaceholder.ListTodo: create request failed: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jsonplaceholder.ListTodo: sending the request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jsonplaceholder.ListTodo: unexpected status code: expected 200; actual %d", resp.StatusCode)
	}

	var todos []Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		return nil, fmt.Errorf("jsonplaceholder.ListTodo: decode response body: %w", err)
	}

	return todos, nil
}

func (c *Client) GetTodo(ctx context.Context, id int64) (*Todo, error) {
	idstr := strconv.FormatInt(id, 10)
	req, err := http.NewRequestWithContext(ctx, "GET", c.listTodoPath+"/"+idstr, nil)
	if err != nil {
		return nil, fmt.Errorf("jsonplaceholder.GetTodo: create request failed: %w", err)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jsonplaceholder.GetTodo: sending the request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrTodoNotFound
		}
		return nil, fmt.Errorf("jsonplaceholder.GetTodo: unexpected status code: expected 200; actual %d", resp.StatusCode)
	}

	var todo Todo
	if err := json.NewDecoder(resp.Body).Decode(&todo); err != nil {
		return nil, fmt.Errorf("jsonplaceholder.GetTodo: decode response body: %w", err)
	}
	return &todo, nil
}
