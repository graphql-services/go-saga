package graphqlorm

import (
	"context"
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/inflection"
	"github.com/machinebox/graphql"
)

// FetchResponse ...
type FetchResponse struct {
	Result interface{} `json:"result"`
}

// GetEntityOptions ...
type GetEntityOptions struct {
	Entity   string
	EntityID *string
	Fields   []string
	Filter   *map[string]interface{}
}

// GetEntity ...
func (c *ORMClient) GetEntity(ctx context.Context, options GetEntityOptions, res interface{}) error {
	query := fmt.Sprintf(`
		query ($id: ID,$filter:%[2]sFilterType) {
			result: %[1]s(id:$id, filter: $filter) {
				id %[3]s
			}
		}
	`, strcase.ToLowerCamel(options.Entity), options.Entity, strings.Join(options.Fields, " "))
	req := graphql.NewRequest(query)
	if options.EntityID != nil {
		req.Var("id", options.EntityID)
	}
	if options.Filter != nil {
		req.Var("filter", options.Filter)
	}

	return c.run(ctx, req, res)
}

// GetEntitiesOptions ...
type GetEntitiesOptions struct {
	Entity string
	Fields []string
	Filter *map[string]interface{}
	Sort   []string
	Limit  *int
	Offset *int
}

// GetEntities ...
func (c *ORMClient) GetEntities(ctx context.Context, options GetEntitiesOptions, res interface{}) error {
	query := fmt.Sprintf(`
		query ($filter: %sFilterType, $sort: [%sSortType!], $limit: Int, $offset: Int) {
			result: %s(filter:$filter,sort:$sort,limit:$limit,offset:$offset) {
				items {
					id %s
				}
				count
			}
		}
	`, options.Entity, options.Entity, inflection.Plural(strcase.ToLowerCamel(options.Entity)), strings.Join(options.Fields, " "))
	req := graphql.NewRequest(query)
	if options.Filter != nil {
		req.Var("filter", options.Filter)
	}
	req.Var("sort", options.Sort)
	if options.Offset != nil {
		req.Var("offset", options.Offset)
	}
	if options.Limit != nil {
		req.Var("limit", options.Limit)
	}

	return c.run(ctx, req, res)
}

// SendQuery ...
func (c *ORMClient) SendQuery(ctx context.Context, query string, variables map[string]interface{}, res interface{}) error {
	req := graphql.NewRequest(query)
	for key, value := range variables {
		req.Var(key, value)
	}
	return c.run(ctx, req, res)
}
