package listener

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/padiazg/notifier-example/config/application"
	"github.com/padiazg/notifier-example/config/settings"
	"github.com/padiazg/notifier-example/listener/amqp"
	"github.com/padiazg/notifier-example/listener/webhook"
)

func Run(settings *settings.Settings) {
	var (
		ctx, cancel              = context.WithCancel(context.Background())
		signalCtx, signalCtxStop = signal.NotifyContext(ctx,
			os.Interrupt,    // interrupt = SIGINT = Ctrl+C
			syscall.SIGQUIT, // Ctrl-\
			syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		)
		app = application.New(settings)
	)

	defer cancel()
	defer signalCtxStop()

	// webhook listener
	if settings.Webhook.Enabled {
		log.Println("Webhook: listener enabled")
		app.Use(webhook.New(&webhook.Config{}))
	}

	// amqp listener
	if settings.AMQP.Enabled {
		log.Println("AMQP: listener enabled")
		app.Use(amqp.New(&amqp.Config{}))
	}

	if len(app.Listeners()) == 0 {
		log.Println("no listeners configured")
		return
	}

	// start listeners
	for _, listener := range app.Listeners() {
		go listener.Run(signalCtx)
	}

	<-signalCtx.Done()

	// shut down listeners
	var wg = &sync.WaitGroup{}
	for _, listener := range app.Listeners() {
		wg.Add(1)
		go listener.Shutdown(wg)
	}

	wg.Wait()
	os.Exit(0)
}
