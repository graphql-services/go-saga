package eventstore

import (
	"context"
	"os"

	"github.com/golang/glog"
	"github.com/machinebox/graphql"
)

func sendRequest(ctx context.Context, req *graphql.Request, data interface{}) error {
	URL := os.Getenv("EVENT_STORE_URL")

	if URL == "" {
		URL = "http://event-store/graphql"
	}

	client := graphql.NewClient(URL)
	client.Log = func(s string) {
		glog.Info(s)
	}

	return client.Run(ctx, req, data)
}
