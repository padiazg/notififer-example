package application

import (
	"context"
	"sync"

	"github.com/padiazg/notifier-example/config/settings"
)

type MicroReceiverInterface interface {
	ConfigureReveiver(*Application)
	Run(context.Context)
	Shutdown(*sync.WaitGroup)
}

type Application struct {
	Settings  *settings.Settings
	listeners []MicroReceiverInterface
}

func New(settings *settings.Settings) *Application {
	return &Application{Settings: settings}
}

func (application *Application) Use(receiver MicroReceiverInterface) {
	receiver.ConfigureReveiver(application)
	application.listeners = append(application.listeners, receiver)
}

func (application *Application) Listeners() []MicroReceiverInterface {
	return application.listeners
}
