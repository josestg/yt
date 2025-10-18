package main

import (
	"cmp"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/josestg/yt/outbound-testing/jsonplaceholder"
)

func main() {
	baseURL := cmp.Or(os.Getenv("JSON_PLACEHOLDER_BASE_URL"), "https://jsonplaceholder.typicode.com")
	c := jsonplaceholder.NewClient(http.DefaultClient, baseURL)
	ctx := context.Background()
	todos, err := c.ListTodo(ctx)
	fmt.Printf("todos: %#v, error: %v", todos, err)
}
