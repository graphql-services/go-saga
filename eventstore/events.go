package eventstore

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
)

// FetchEventsSort ...
type FetchEventsSort string

// FetchEventType ...
type FetchEventType string

const (
	// FetchEventsSortDateAsc ...
	FetchEventsSortDateAsc FetchEventsSort = "DATE_ASC"
	// FetchEventsSortDateDesc ...
	FetchEventsSortDateDesc FetchEventsSort = "DATE_DESC"

	// FetchEventTypeCreated ...
	FetchEventTypeCreated = "CREATED"
	// FetchEventTypeUpdated ...
	FetchEventTypeUpdated = "UPDATED"
	// FetchEventTypeDeleted ...
	FetchEventTypeDeleted = "DELETED"
)

// FetchEventsOptions ...
type FetchEventsOptions struct {
	ID     *string
	Entity *string
	Cursor *string
	Limit  *int
	Sort   *FetchEventsSort
}

// FetchEventsResponse ...
type FetchEventsResponse struct {
	Events []struct {
		ID            string `json:"id"`
		Cursor        string
		OperationName string
		Type          FetchEventType
		Entity        string
		EntityID      string `json:"entityId"`
		Data          string
		Date          time.Time
		Columns       []string `json:"columns"`
		PrincipalID   *string  `json:"principalId"`
	}
}

// FetchEvents ...
func FetchEvents(ctx context.Context, options FetchEventsOptions, data *FetchEventsResponse) error {
	query := `
		query ($id:ID, $cursor: String, $limit: Int = 100, $entity: EventEntities, $sort: EventEntitiesSort) {
			events: _events(id:$id,cursor:$cursor,limit:$limit,entity: $entity, sort: $sort) {
				id
				cursor
				operationName
				type
				entity
				entityId
				data
				date
				columns
				principalId
			}
		}
	`
	req := graphql.NewRequest(query)

	if options.ID != nil {
		req.Var("id", *options.ID)
	}
	if options.Entity != nil {
		req.Var("entity", *options.Entity)
	}
	if options.Cursor != nil {
		req.Var("cursor", *options.Cursor)
	}
	if options.Limit != nil {
		req.Var("limit", *options.Limit)
	}
	if options.Sort != nil {
		req.Var("sort", *options.Sort)
	}

	return sendRequest(ctx, req, data)
}
