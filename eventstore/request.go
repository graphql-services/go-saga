package eventstore

import (
	"context"
	"fmt"
	"os"

	"github.com/golang/glog"
	"github.com/machinebox/graphql"
)

func sendRequest(ctx context.Context, req *graphql.Request, data interface{}) error {
	URL := os.Getenv("EVENT_STORE_URL")

	if URL == "" {
		return fmt.Errorf("Missing required environment variable EVENT_STORE_URL")
	}

	client := graphql.NewClient(URL)
	client.Log = func(s string) {
		glog.Info(s)
	}

	return client.Run(ctx, req, data)
}
