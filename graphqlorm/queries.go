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
	EntityID string
	Fields   []string
	Filter   map[string]interface{}
}

// GetEntity ...
func (c *ORMClient) GetEntity(ctx context.Context, options GetEntityOptions, res interface{}) error {
	query := fmt.Sprintf(`
		query ($id: ID,$filter:%[1]sFilterType) {
			result: %[1]s(id:$id, filter: $filter) {
				id %[2]s
			}
		}
	`, strcase.ToLowerCamel(options.Entity), strings.Join(options.Fields, " "))
	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)
	req.Var("filter", options.Filter)

	return c.run(ctx, req, res)
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
	req.Var("filter", options.Filter)
	req.Var("sort", options.Sort)
	if options.Offset != nil {
		req.Var("offset", options.Offset)
	}
	if options.Limit != nil {
		req.Var("limit", options.Limit)
	}

	return c.run(ctx, req, res)
}
