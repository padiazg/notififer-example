package emitter

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/padiazg/notifier-example/config/settings"
	"github.com/padiazg/notifier/engine"
	"github.com/padiazg/notifier/model"

	amqp "github.com/padiazg/notifier/connector/amqp10"
	"github.com/padiazg/notifier/connector/webhook"
)

const (
	SomeEvent    model.EventType = "SomeEvent"
	AnotherEvent model.EventType = "AnotherEvent"
)

func Run(settings *settings.Settings) {
	var (
		engine = engine.New(&engine.Config{
			OnError: func(err error) {
				log.Printf("Error: %s", err.Error())
			},
		})

		webhookId, amqpId string

		wg   sync.WaitGroup
		done = make(chan bool)
		c    = make(chan os.Signal, 2)

		messages = []*model.Notification{
			{
				Event: SomeEvent,
				Data:  "simple text data",
			},
			{
				Event: AnotherEvent,
				Data: struct {
					ID   uint
					Name string
				}{ID: 1, Name: "complex data"},
			},
			{
				Event:    SomeEvent,
				Data:     "only to webhook",
				Channels: []string{webhookId},
			},
			{
				Event:    SomeEvent,
				Data:     "only to mq",
				Channels: []string{amqpId},
			},
		}
	)

	if settings.Webhook.Enabled {
		webhookId = engine.Register(webhook.New(&webhook.Config{
			Name:     "Webhook",
			Endpoint: "https://localhost:4443/webhook",
			Insecure: true,
			Headers: map[string]string{
				"Authorization": "Bearer xyz123",
				"X-Portal-Id":   "1234567890",
			},
		}))
	}

	if settings.AMQP.Enabled {
		amqpId = engine.Register(amqp.New(&amqp.Config{
			Name:      "AMQP",
			QueueName: "notifier",
			Address:   "amqp://localhost",
		}))
	}

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Starting engine...")
	engine.Start()

	for _, m := range messages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			engine.Dispatch(m)
		}()
	}

	wg.Wait()

	// close the done channel and exit the program
	close(done)

	for {
		select {
		case <-done:
			fmt.Println("Exiting...")
			// stop the engine
			engine.Stop()
			// give some time to the engine to stop
			time.Sleep(100 * time.Millisecond)
			return

		case s := <-c:
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Got SIGINT/SIGTERM, shutting down emmiter...")
				close(done)
			}
		}
	}
}
