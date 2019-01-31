package eventstore

import (
	"context"
	"time"

	"github.com/machinebox/graphql"
)

// FetchEventsSort ...
type FetchEventsSort string

const (
	// FetchEventsSortDateAsc ...
	FetchEventsSortDateAsc FetchEventsSort = "DATE_ASC"
	// FetchEventsSortDateDesc ...
	FetchEventsSortDateDesc FetchEventsSort = "DATE_DESC"
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
		ID            string         `json:"id"`
		Cursor        string         `json:"cursor"`
		OperationName *string        `json:"operationName"`
		Entity        string         `json:"entity"`
		EntityID      string         `json:"entityId"`
		Data          string         `json:"data"`
		OldValues     []EventValue   `json:"oldValues"`
		NewValues     []EventValue   `json:"newValues"`
		Type          FetchEventType `json:"type"`
		Date          time.Time      `json:"date"`
		PrincipalID   *string        `json:"principalId"`
		Columns       []string       `json:"columns"`
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
				newValues { name value }
				oldValues { name value }
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
