package aggregator

import (
	"context"
	"fmt"
	"strings"

	"os"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
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

// GetEntitiesOptions ...
type GetEntitiesOptions struct {
	Entity string
	Fields []string
	Filter map[string]interface{}
	Sort   []string
	Limit  *int
	Offset *int
}

// GetEntities ...
func GetEntities(options GetEntitiesOptions, res interface{}) error {
	query := fmt.Sprintf(`
		query ($filter: %sFilterType, $sort: [%sSortType], $limit: Int, $offset: Int) {
			result: %s(filter:$filter,sort:$sort,limit:$limit,offset:$offset) {
				items {
					id %s
				}
				count
			}
		}
	`, options.Entity, options.Entity, inflection.Plural(strcase.ToLowerCamel(options.Entity)), strings.Join(options.Fields, " "))
	req := graphql.NewRequest(query)
	req.Var("filter", options.Filter)
	req.Var("sort", options.Sort)
	if options.Offset != nil {
		req.Var("offset", options.Offset)
	}
	if options.Limit != nil {
		req.Var("limit", options.Limit)
	}

	return sendRequest(req, res)
}
