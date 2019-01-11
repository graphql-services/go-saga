package aggregator

import (
	"context"
	"fmt"
	"strings"

	"os"

	"github.com/iancoleman/strcase"
	"github.com/machinebox/graphql"
)

// FetchResponse ...
type FetchResponse struct {
	Result interface{} `json:"result"`
}

func sendRequest(req *graphql.Request, res interface{}) error {
	ctx := context.Background()
	URL := os.Getenv("AGGREGATOR_URL")

	if URL == "" {
		URL = "http://event-store-aggergator/graphql"
	}

	client := graphql.NewClient(URL)
	client.Log = func(s string) {
		fmt.Println(s)
	}

	return client.Run(ctx, req, &res)
}

// GetEntityOptions ...
type GetEntityOptions struct {
	Entity   string
	EntityID string
	Fields   []string
}

// GetEntity ...
func GetEntity(options GetEntityOptions, res interface{}) error {
	query := fmt.Sprintf(`
		query ($id: ID!) {
			result: %s(id:$id) {
				id %s
			}
		}
	`, strcase.ToLowerCamel(options.Entity), strings.Join(options.Fields, " "))
	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)

	return sendRequest(req, res)
}
