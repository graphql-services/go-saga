package eventstore

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bitly/go-nsq"
)

// Event ...
type Event struct {
	ID          string      `json:"id"`
	Entity      string      `json:"entity"`
	EntityID    string      `json:"entityId"`
	Data        interface{} `json:"data"`
	Type        string      `json:"type"`
	Date        time.Time   `json:"date"`
	PrincipalID *string     `json:"principalId"`
	Columns     []string    `json:"columns"`
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

// OnEventOptions ...
type OnEventOptions struct {
	Channel     string
	HandlerFunc func(e Event) error
}

// OnEvent ...
func OnEvent(options OnEventOptions) {
	URL := os.Getenv("NSQ_URL")
	LOOKUP_URL := os.Getenv("NSQ_LOOKUP_URL")

	if URL == "" && LOOKUP_URL == "" {
		log.Panic("You have to specify NSQ_URL or NSQ_LOOKUP_URL in environment variables")
	}

	go func() {
		config := nsq.NewConfig()
		q, err := nsq.NewConsumer("es-event", options.Channel, config)
		if err != nil {
			log.Panic(err)
		}
		q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
			var e Event
			err := json.Unmarshal(message.Body, &e)
			if err != nil {
				return err
			}
			return options.HandlerFunc(e)
		}))

		if URL != "" {
			addresses := strings.Split(URL, ",")
			err = q.ConnectToNSQDs(addresses)
			if err != nil {
				log.Panic(err)
			}
		}

		if LOOKUP_URL != "" {
			addresses := strings.Split(LOOKUP_URL, ",")
			err = q.ConnectToNSQLookupds(addresses)
			if err != nil {
				log.Panic(err)
			}
		}
	}()
}
