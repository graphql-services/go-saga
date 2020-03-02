package graphqlorm

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/machinebox/graphql"
)

// ORMClient ...
type ORMClient struct {
	gc *graphql.Client
}

// NewClient ...
func NewClient(URL string) *ORMClient {
	client := graphql.NewClient(URL, graphql.WithHTTPClient(xray.Client(http.DefaultClient)))
	if os.Getenv("DEBUG") == "true" {
		client.Log = func(s string) { log.Println(s) }
	}
	return &ORMClient{client}
}

func (c *ORMClient) run(ctx context.Context, req *graphql.Request, data interface{}) error {
	return c.gc.Run(ctx, req, data)
}
