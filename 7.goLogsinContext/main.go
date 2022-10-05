package main

import (
	"log"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigInfoLogger(os.Stdout),

		// Workshop > Let the agent collect, and forward logs automatically
		// https://docs.newrelic.com/docs/logs/logs-context/configure-logs-context-go/
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		panic(err)
	}

	app.WaitForConnection(5 * time.Second)

	// Create a logWriter, then pass it to the log.Logger
	writer := logWriter.New(os.Stdout, app)
	logger := log.New(&writer, "Background:  ", log.Default().Flags())

	logger.Print("Hello world!")

	txnName := "logsInContext Sample Transaction"
	txn := app.StartTransaction(txnName)

	// Always create a new log object in order to avoid changing the context of the logger for
	// other threads that may be logging outside of this transaction
	txnLogger := log.New(writer.WithTransaction(txn), "Transaction: ", log.Default().Flags())
	txnLogger.Printf("In transaction %s.", txnName)

	// simulate doing something
	time.Sleep(time.Microsecond * 100)

	txnLogger.Printf("Ending transaction %s.", txnName)
	txn.End()

	logger.Print("Goodbye!")
	app.Shutdown(10 * time.Second)
}
