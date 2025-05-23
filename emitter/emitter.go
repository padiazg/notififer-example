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
	"github.com/padiazg/notifier/broker"
	amqp "github.com/padiazg/notifier/connector/amqp10"
	"github.com/padiazg/notifier/connector/webhook"
	"github.com/padiazg/notifier/model"
)

const (
	SomeEvent    model.EventType = "SomeEvent"
	AnotherEvent model.EventType = "AnotherEvent"
)

func Run(settings *settings.Settings) {
	var (
		broker = broker.New(&broker.Config{
			OnError: func(err error) {
				log.Printf("Error: %s", err.Error())
			},
		})

		webhookId string
		amqpId    string

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
		}
	)

	if settings.Webhook.Enabled {
		var (
			webhookConfig = &webhook.Config{
				Name: "Webhook",
				Headers: map[string]string{
					"Authorization": "Bearer xyz123",
					"X-Portal-Id":   "1234567890",
				},
			}
			proto string = "http"
		)

		if settings.Webhook.UseTLS {
			proto = "https"
			webhookConfig.Insecure = true
		}

		webhookConfig.Endpoint = fmt.Sprintf("%s://localhost:%d/webhook", proto, settings.Webhook.Port)
		webhookId = broker.NotifierRegister(webhook.New(webhookConfig))

		messages = append(messages, &model.Notification{
			Event:    SomeEvent,
			Data:     "only to webhook",
			Channels: []string{webhookId},
		})
	}

	if settings.AMQP.Enabled {
		amqpId = broker.NotifierRegister(amqp.New(&amqp.Config{
			Name:      "AMQP",
			QueueName: settings.AMQP.Queue,
			Address:   settings.AMQP.Address,
		}))

		messages = append(messages, &model.Notification{
			Event:    SomeEvent,
			Data:     "only to mq",
			Channels: []string{amqpId},
		})
	}

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Starting broker...")
	broker.Start()

	for _, m := range messages {
		wg.Add(1)
		go func() {
			defer wg.Done()
			broker.Dispatch(m)
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
			broker.Stop()
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
