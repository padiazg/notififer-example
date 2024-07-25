package webhook

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/padiazg/notifier-example/config/application"
)

type Config struct {
}

type Webhook struct {
	*Config
	mux    *http.ServeMux
	server *http.Server
}

func New(config *Config) *Webhook {
	return &Webhook{Config: config}
}

func (r *Webhook) ConfigureReveiver(application *application.Application) {
	r.mux = http.NewServeMux()
	r.mux.HandleFunc("/webhook", handleWebhook)

	r.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", application.Settings.Webhook.Port),
		Handler: r.mux,
	}
}

func (r *Webhook) Run(ctx context.Context) {
	log.Printf("Webhook: server listening on %s\n", r.server.Addr)

	r.server.BaseContext = func(_ net.Listener) context.Context { return ctx }

	if err := r.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Webhook: Starting server: %+v", err)
	}
}

func (r *Webhook) Shutdown(wg *sync.WaitGroup) {
	log.Printf("Webhook: Shutting down\n")
	defer wg.Done()

	const timeout = 5 * time.Second
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()

	if err := r.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Webhook: %+v", err)
	}
}
