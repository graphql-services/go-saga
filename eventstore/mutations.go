package eventstore

import (
	"context"
	"fmt"

	"os"

	"github.com/machinebox/graphql"
)

func sendRequest(req *graphql.Request) (interface{}, error) {
	ctx := context.Background()
	URL := os.Getenv("EVENT_STORE_URL")

	if URL == "" {
		URL = "http://event-store/graphql"
	}

	client := graphql.NewClient(URL)

	var data interface{}
	err := client.Run(ctx, req, &data)
	return data, err
}

// CreateEntityOptions ...
type CreateEntityOptions struct {
	Entity string      `json:"entity"`
	Input  interface{} `json:"input"`
}

// CreateEntity ...
func CreateEntity(options CreateEntityOptions) (interface{}, error) {
	query := fmt.Sprintf(`
		mutation ($input: %sRawCreateInput!) {
			create%s (input:$input) {
				id
			}
		}
	`, options.Entity, options.Entity)
	req := graphql.NewRequest(query)
	req.Var("input", options.Input)

	return sendRequest(req)
}

// UpdateEntityOptions ...
type UpdateEntityOptions struct {
	Entity   string      `json:"entity"`
	EntityID string      `json:"entityId"`
	Input    interface{} `json:"input"`
}

// UpdateEntity ...
func UpdateEntity(options UpdateEntityOptions) (interface{}, error) {
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

	return sendRequest(req)
}

// DeleteEntityOptions ...
type DeleteEntityOptions struct {
	Entity   string `json:"entity"`
	EntityID string `json:"entityId"`
}

// DeleteEntity ...
func DeleteEntity(options DeleteEntityOptions) (interface{}, error) {
	query := fmt.Sprintf(`
		mutation ($id:ID!) {
			delete%s (id:$id) {
				id
			}
		}
	`, options.Entity)

	req := graphql.NewRequest(query)
	req.Var("id", options.EntityID)

	return sendRequest(req)
}
