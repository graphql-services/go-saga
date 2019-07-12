package graphqlorm

import (
	"context"
	"fmt"

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
func (c *ORMClient) CreateEntity(ctx context.Context, options CreateEntityOptions) (MutationResult, error) {
	query := fmt.Sprintf(`
		mutation ($input: %sCreateInput!) {
			result: create%s (input:$input) {
				id
			}
		}
	`, options.Entity, options.Entity)
	req := graphql.NewRequest(query)
	req.Var("input", options.Input)

	var data struct {
		Result MutationResult
	}
	err := c.run(ctx, req, &data)
	return data.Result, err
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
