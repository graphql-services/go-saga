package graphqlorm

import (
	"context"
	"fmt"
	"strings"

	"github.com/machinebox/graphql"
)

// CreateEntityOptions ...
type CreateEntityOptions struct {
	Entity string      `json:"entity"`
	Input  interface{} `json:"input"`
}

// MutationResult ...
type MutationResult struct {
	ID string `json:"id"`
}

// CreateEntity ...
func (c *ORMClient) CreateEntity(ctx context.Context, options CreateEntityOptions) (res MutationResult, err error) {
	_res, err := c.CreateEntities(ctx, []CreateEntityOptions{options})
	if err == nil {
		res = _res[0]
	}
	return
}

// CreateEntities ...
func (c *ORMClient) CreateEntities(ctx context.Context, options []CreateEntityOptions) ([]MutationResult, error) {

	inputs := []string{}
	for key, value := range options {
		inputs = append(inputs, fmt.Sprintf(`$input%d: %sCreateInput!`, key, value.Entity))
	}
	results := []string{}
	for key, value := range options {
		results = append(results, fmt.Sprintf(`result%d: create%s (input:$input%d) {
			id
		}
		`, key, value.Entity, key))
	}

	query := fmt.Sprintf(`
		mutation (%s) {
			%s
		}
	`, strings.Join(inputs, ","), strings.Join(results, "\n"))
	fmt.Println("running query", query)
	req := graphql.NewRequest(query)
	for key, value := range options {
		req.Var(fmt.Sprintf("input%d", key), value.Input)
	}

	var data map[string]MutationResult
	err := c.run(ctx, req, &data)

	res := []MutationResult{}

	for _, val := range data {
		res = append(res, val)
	}

	return res, err
}

// UpdateEntityOptions ...
type UpdateEntityOptions struct {
	Entity   string      `json:"entity"`
	EntityID string      `json:"entityId"`
	Input    interface{} `json:"input"`
}

// UpdateEntity ...
func (c *ORMClient) UpdateEntity(ctx context.Context, options UpdateEntityOptions) (MutationResult, error) {
	query := fmt.Sprintf(`
		mutation ($id:ID!, $input: %sUpdateInput!) {
			update%s (id:$id, input:$input) {
				id
			}
		}
	`, options.Entity, options.Entity)

	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)
	req.Var("input", options.Input)

	var data struct {
		Result MutationResult
	}
	err := c.run(ctx, req, &data)
	return data.Result, err
}

// DeleteEntityOptions ...
type DeleteEntityOptions struct {
	Entity   string `json:"entity"`
	EntityID string `json:"entityId"`
}

// DeleteEntity ...
func (c *ORMClient) DeleteEntity(ctx context.Context, options DeleteEntityOptions) (MutationResult, error) {
	query := fmt.Sprintf(`
		mutation ($id:ID!) {
			delete%s (id:$id) {
				id
			}
		}
	`, options.Entity)

	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)

	var data struct {
		Result MutationResult
	}
	err := c.run(ctx, req, &data)
	return data.Result, err
}
