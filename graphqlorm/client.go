package graphqlorm

import (
	"context"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

// ORMClient ...
type ORMClient struct {
	gc *graphql.Client
}

// NewClient ...
func NewClient(URL string) *ORMClient {
	client := graphql.NewClient(URL)
	if os.Getenv("DEBUG") == "true" {
		client.Log = func(s string) { log.Println(s) }
	}
	return &ORMClient{client}
}

func (c *ORMClient) run(ctx context.Context, req *graphql.Request, data interface{}) error {
	return c.gc.Run(ctx, req, data)
}
