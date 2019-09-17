package graphqlorm

import (
	"context"
	"log"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/novacloudcz/graphql-orm/events"
)

type EventHandler (func(event events.Event) error)

type HandleEventOptions struct {
	port string
}

func HandleEvent(handler EventHandler, opts *HandleEventOptions) {
	portString := opts.port
	if portString == "" {
		portString = "80"
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic(err)
	}

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithPort(port),
		cloudevents.WithStructuredEncoding(),
	)
	if err != nil {
		panic(err)
	}
	c, err := cloudevents.NewClient(t)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Printf("listening on http://localhost:%d", port)
	log.Fatal(c.StartReceiver(context.Background(), func(event cloudevents.Event) error {
		var ormEvent events.Event
		err := event.DataAs(&ormEvent)
		if err != nil {
			return err
		}
		return handler(ormEvent)
	}))
}
