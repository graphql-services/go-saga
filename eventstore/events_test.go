package eventstore

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/matryer/is"
)

func TestFetch(t *testing.T) {
	is := is.New(t)

	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Header().Add("content-type", "application/json")
		io.WriteString(w, `{
			"data":{
				"events": [{
					"id": "blah",
					"entity": "User",
					"entityId": "john-doe-id",
					"cursor":"abcd1",
					"operationName":"CreateUser",
					"type":"CREATED",
					"data":"{\"username\":\"john.doe\"}",
					"date":"2018-01-02T06:23:00Z",
					"columns":["username"],
					"principalId":"administrator-id"
				}]
			}
		}`)
	}))
	os.Setenv("EVENT_STORE_URL", srv.URL)
	defer srv.Close()

	ctx := context.Background()
	options := FetchEventsOptions{}
	var res FetchEventsResponse
	err := FetchEvents(ctx, options, &res)

	is.NoErr(err)
	is.Equal(len(res.Events), 1)
	is.Equal(len(res.Events[0].Columns), 1)
	is.Equal(res.Events[0].Columns[0], "username")
	is.Equal(res.Events[0].Date.Day(), 2)
	is.Equal(res.Events[0].Date.Year(), 2018)

	is.Equal(calls, 1)
	// is.Equal(responseData["something"], "yes")
}
