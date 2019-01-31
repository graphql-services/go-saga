package eventstore

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	vice "github.com/matryer/vice/queues/nsq"
	"github.com/nsqio/go-nsq"
)

// FetchEventType ...
type FetchEventType string

const (
	// FetchEventTypeCreated ...
	FetchEventTypeCreated = "CREATED"
	// FetchEventTypeUpdated ...
	FetchEventTypeUpdated = "UPDATED"
	// FetchEventTypeDeleted ...
	FetchEventTypeDeleted = "DELETED"
)

// EventValue ...
type EventValue struct {
	Name  string  `json:"name"`
	Value *string `json:"value"`
}

// Event ...
type Event struct {
	ID            string         `json:"id"`
	Cursor        string         `json:"cursor"`
	OperationName *string        `json:"operationName"`
	Entity        string         `json:"entity"`
	EntityID      string         `json:"entityId"`
	Data          interface{}    `json:"data"`
	OldValues     []EventValue   `json:"oldValues"`
	NewValues     []EventValue   `json:"newValues"`
	Type          FetchEventType `json:"type"`
	Date          time.Time      `json:"date"`
	PrincipalID   *string        `json:"principalId"`
	Columns       []string       `json:"columns"`
}

// HasColumn check if given event has changes on specific column
func (e Event) HasColumn(c string) bool {
	for _, col := range e.Columns {
		if col == c {
			return true
		}
	}
	return false
}

// OldValue ...
func (e Event) OldValue(c string) *EventValue {
	for _, v := range e.OldValues {
		if v.Name == c {
			return &v
		}
	}
	return nil
}

// NewValue ...
func (e Event) NewValue(c string) *EventValue {
	for _, v := range e.NewValues {
		if v.Name == c {
			return &v
		}
	}
	return nil
}

// OnEventOptions ...
type OnEventOptions struct {
	Topic       string
	Channel     string
	HandlerFunc func(e Event) error
}

// OnEvent ...
func OnEvent(options OnEventOptions) {
	URL := os.Getenv("NSQ_URL")
	lookupURL := os.Getenv("NSQ_LOOKUP_URL")

	if URL == "" && lookupURL == "" {
		log.Panic("You have to specify NSQ_URL or NSQ_LOOKUP_URL in environment variables")
	}

	go func() {
		transport := vice.New()
		transport.NewConsumer = func(name string) (*nsq.Consumer, error) {
			return nsq.NewConsumer(name, options.Channel, nsq.NewConfig())
		}
		transport.ConnectConsumer = func(consumer *nsq.Consumer) error {
			if URL != "" {
				addresses := strings.Split(URL, ",")
				fmt.Println("connect NSQ", addresses)
				return consumer.ConnectToNSQDs(addresses)
			}
			addresses := strings.Split(lookupURL, ",")
			return consumer.ConnectToNSQLookupds(addresses)
		}

		topic := options.Topic
		if topic == "" {
			topic = "es-event"
		}

		events := transport.Receive(topic)
		errChan := transport.ErrChan()

		for {
			select {
			case event := <-events:
				var e Event
				err := json.Unmarshal(event, &e)
				if err != nil {
					panic(err)
				}
				if err := options.HandlerFunc(e); err != nil {
					panic(err)
				}
			case err := <-errChan:
				panic(err)
			}
		}
	}()
}
