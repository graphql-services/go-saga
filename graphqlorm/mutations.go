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
func (c *ORMClient) UpdateEntity(ctx context.Context, options UpdateEntityOptions) (res MutationResult, err error) {
	_res, err := c.UpdateEntities(ctx, []UpdateEntityOptions{options})
	if err == nil {
		res = _res[0]
	}
	return
}

// UpdateEntities ...
func (c *ORMClient) UpdateEntities(ctx context.Context, options []UpdateEntityOptions) ([]MutationResult, error) {

	inputs := []string{}
	for key, value := range options {
		inputs = append(inputs, fmt.Sprintf(`$id%d: ID!, $input%d: %sUpdateInput!`, key, key, value.Entity))
	}
	results := []string{}
	for key, value := range options {
		results = append(results, fmt.Sprintf(`result%d: update%s (id:$id%d, input:$input%d) {
			id
		}
		`, key, value.Entity, key, key))
	}

	query := fmt.Sprintf(`
		mutation (%s) {
			%s
		}
	`, strings.Join(inputs, ","), strings.Join(results, "\n"))
	req := graphql.NewRequest(query)
	for key, value := range options {
		req.Var(fmt.Sprintf("id%d", key), value.EntityID)
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

// DeleteEntityOptions ...
type DeleteEntityOptions struct {
	Entity   string `json:"entity"`
	EntityID string `json:"entityId"`
}

// DeleteEntity ...
func (c *ORMClient) DeleteEntity(ctx context.Context, options DeleteEntityOptions) (res MutationResult, err error) {
	_res, err := c.DeleteEntities(ctx, []DeleteEntityOptions{options})
	if err == nil {
		res = _res[0]
	}
	return
}

// DeleteEntities ...
func (c *ORMClient) DeleteEntities(ctx context.Context, options []DeleteEntityOptions) ([]MutationResult, error) {

	inputs := []string{}
	for key := range options {
		inputs = append(inputs, fmt.Sprintf(`$id%d: ID!`, key))
	}
	results := []string{}
	for key, value := range options {
		results = append(results, fmt.Sprintf(`result%d: delete%s (id:$id%d) {
			id
		}
		`, key, value.Entity, key))
	}

	query := fmt.Sprintf(`
		mutation (%s) {
			%s
		}
	`, strings.Join(inputs, ","), strings.Join(results, "\n"))

	req := graphql.NewRequest(query)
	for key, value := range options {
		req.Var(fmt.Sprintf("id%d", key), value.EntityID)
	}

	var data map[string]MutationResult
	err := c.run(ctx, req, &data)

	res := []MutationResult{}

	for _, val := range data {
		res = append(res, val)
	}

	return res, err
}
