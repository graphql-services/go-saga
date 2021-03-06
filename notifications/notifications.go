package notifications

import (
	"context"
	"fmt"
	"time"

	"os"

	"github.com/machinebox/graphql"
)

// DEPRECATED, please use github.com/graphql-services/notifications/client

func sendRequest(req *graphql.Request) (interface{}, error) {
	ctx := context.Background()
	URL := os.Getenv("NOTIFICATIONS_URL")

	if URL == "" {
		return nil, fmt.Errorf("Missing required environment variable NOTIFICATIONS_URL")
	}

	client := graphql.NewClient(URL)
	client.Log = func(s string) {
		fmt.Println(s)
	}

	var data interface{}
	err := client.Run(ctx, req, &data)
	return data, err
}

// CreateNotificationInput ...
type CreateNotificationInput struct {
	Message     string    `json:"message"`
	Principal   *string   `json:"principal"`
	Channel     *string   `json:"channel"`
	Reference   *string   `json:"reference"`
	ReferenceID *string   `json:"referenceID"`
	Date        time.Time `json:"date"`
}

// CreateNotification ...
func CreateNotification(input CreateNotificationInput) (interface{}, error) {
	query := fmt.Sprintf(`
		mutation ($input: EventStoreNotificationInput!) {
			createNotification (input:$input) {
				id
				message
				date
				principal
				channel
				reference
				referenceID
			}
		}
	`)
	req := graphql.NewRequest(query)
	req.Var("input", input)

	return sendRequest(req)
}
