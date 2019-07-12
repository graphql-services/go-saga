package graphqlorm

import (
	"context"

	"github.com/machinebox/graphql"
)

// ORMClient ...
type ORMClient struct {
	gc *graphql.Client
}

// NewClient ...
func NewClient(URL string) *ORMClient {
	client := graphql.NewClient(URL)
	return &ORMClient{client}
}

func (c *ORMClient) run(ctx context.Context, req *graphql.Request, data interface{}) error {
	return c.gc.Run(ctx, req, data)
}
