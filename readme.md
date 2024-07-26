# notififer-example
An example to demonstrate how to use the [notifier](https://github.com/padiazg/notifier) library

## Build
```bash
$ go build -o bin/notifier-example
# using make
$ make build
```
 
## Usage
The application has two components: an emitter and a listener. 

### Configuration file
The default configuration file is `.notifier-example.yaml` and it's commmon to both components. You can create an example configuration file by running `bin/notifier-example create`
 
### Listener
The listener will start all the listeners set in the configuration file and waits for notifications. You can stop it by hitting `CTRL+C` or `CTRL+\`
```bash
$ bin/notifier-example listen
## using make
$ make listen
```

### Emitter
The emitter will send a few messages to the listeners enabled in the configuration file and then stop immediately
```bash
$ bin/notifier-example emmit
## using make
$ make emmit
```
