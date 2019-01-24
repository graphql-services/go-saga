package eventstore

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

// CreateEntity ...
func CreateEntity(ctx context.Context, options CreateEntityOptions) (interface{}, error) {
	query := fmt.Sprintf(`
		mutation ($input: %sRawCreateInput!) {
			create%s (input:$input) {
				id
			}
		}
	`, options.Entity, options.Entity)
	req := graphql.NewRequest(query)
	req.Var("input", options.Input)

	var data interface{}
	err := sendRequest(ctx, req, &data)
	return data, err
}

// UpdateEntityOptions ...
type UpdateEntityOptions struct {
	Entity   string      `json:"entity"`
	EntityID string      `json:"entityId"`
	Input    interface{} `json:"input"`
}

// UpdateEntity ...
func UpdateEntity(ctx context.Context, options UpdateEntityOptions) (interface{}, error) {
	query := fmt.Sprintf(`
		mutation ($id:ID!, $input: %sRawUpdateInput!) {
			update%s (id:$id, input:$input) {
				id
			}
		}
	`, options.Entity, options.Entity)

	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)
	req.Var("input", options.Input)

	var data interface{}
	err := sendRequest(ctx, req, &data)
	return data, err
}

// DeleteEntityOptions ...
type DeleteEntityOptions struct {
	Entity   string `json:"entity"`
	EntityID string `json:"entityId"`
}

// DeleteEntity ...
func DeleteEntity(ctx context.Context, options DeleteEntityOptions) (interface{}, error) {
	query := fmt.Sprintf(`
		mutation ($id:ID!) {
			delete%s (id:$id) {
				id
			}
		}
	`, options.Entity)

	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)

	var data interface{}
	err := sendRequest(ctx, req, &data)
	return data, err
}
