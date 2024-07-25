package amqp

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/padiazg/notifier-example/config/application"
	m "github.com/padiazg/notifier/model"
)

type Config struct {
}

type AMQP struct {
	*Config
	conn     *amqp.Conn
	session  *amqp.Session
	receiver *amqp.Receiver
	ctx      context.Context
	cancel   context.CancelFunc
}

func New(config *Config) *AMQP {
	log.Printf("AMQP: New")
	return &AMQP{Config: config}
}

func (r *AMQP) ConfigureReveiver(application *application.Application) {
	var err error

	log.Printf("AMQP: ConfigureReveiver")
	if r.conn, err = amqp.Dial(context.TODO(), "amqp:localhost", nil); err != nil {
		log.Fatal("Dialing AMQP server:", err)
	}

	if r.session, err = r.conn.NewSession(context.TODO(), nil); err != nil {
		log.Fatal("Creating AMQP session:", err)
	}
}

func (r *AMQP) Run(ctx context.Context) {
	var err error

	r.ctx, r.cancel = context.WithCancel(ctx)

	log.Printf("AMQP: server started\n")

	// {
	// create a receiver
	r.receiver, err = r.session.NewReceiver(r.ctx, "notifier", &amqp.ReceiverOptions{
		SourceDurability: amqp.DurabilityUnsettledState,
	})
	if err != nil {
		log.Fatal("Creating receiver link:", err)
	}

	log.Printf("Connected to queue\n")

	for {
		select {
		case <-r.ctx.Done():
			log.Printf("AMQP: done received at Run for")
			return

		default:
			// receive next message
			msg, err := r.receiver.Receive(r.ctx, nil)
			if err != nil {
				log.Printf("Reading message from AMQP: %v\n", err)
			}

			// accept message
			if err = r.receiver.AcceptMessage(r.ctx, msg); err != nil {
				log.Printf("Failure accepting message: %v", err)
			}

			r.showMessage(msg.GetData())
		}

	}
}

func (r *AMQP) Shutdown(wg *sync.WaitGroup) {
	log.Printf("AMQP: Shutting down\n")
	defer wg.Done()

	// stop Run loop
	r.cancel()

	log.Printf("AMQP: Releasing receiver\n")
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 1*time.Second)
	r.receiver.Close(shutdownCtx)
	defer shutdownRelease()

	log.Printf("AMQP: Closing connection\n")
	if err := r.conn.Close(); err != nil {
		log.Fatalf("AMQP: %+v", err)
	}
}

func (r *AMQP) showMessage(msg []byte) {
	// Parse the incoming JSON payload
	var notification m.Notification

	if err := json.Unmarshal(msg, &notification); err != nil {
		log.Printf("Failed to decode JSON payload: %v", err)
		return
	}

	formated, err := json.MarshalIndent(notification, "", "  ")
	if err != nil {
		log.Printf("Failed to format payload: %v", err)
		return
	}

	log.Printf("Message received: %s\n", formated)
}
