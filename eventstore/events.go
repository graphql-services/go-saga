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
	ID         *string
	Entity     *string
	CursorFrom *string
	Limit      *int
	Sort       *FetchEventsSort
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
		OldValues     *map[string]interface{} `json:"oldValues"`
		NewValues     *map[string]interface{} `json:"newValues"`
		Columns       []string                `json:"columns"`
		PrincipalID   *string                 `json:"principalId"`
	}
}

// FetchEvents ...
func FetchEvents(ctx context.Context, options FetchEventsOptions, data *FetchEventsResponse) error {
	query := `
		query ($id:ID, $cursorFrom: String, $limit: Int = 100, $entity: EventEntities, $sort: EventEntitiesSort) {
			events: _events(id:$id,cursorFrom:$cursorFrom,limit:$limit,entity: $entity, sort: $sort) {
				id
				cursor
				operationName
				type
				entity
				entityId
				data
				newValues
				oldValues
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
	if options.CursorFrom != nil {
		req.Var("cursorFrom", *options.CursorFrom)
	}
	if options.Limit != nil {
		req.Var("limit", *options.Limit)
	}
	if options.Sort != nil {
		req.Var("sort", *options.Sort)
	}

	return sendRequest(ctx, req, data)
}
