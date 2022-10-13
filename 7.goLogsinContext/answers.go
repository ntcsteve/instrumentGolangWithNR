package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/logWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// init sets initial values for variables used in the function.
func init() {
	rand.Seed(time.Now().UnixNano())
}

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

	// Not necessary for monitoring a production application with a lot of data.
	app.WaitForConnection(5 * time.Second)

	// Workshop - Create a logWriter, then pass it to the log.Logger
	// https://docs.newrelic.com/docs/logs/logs-context/configure-logs-context-go/#1-standard-library-log
	writer := logWriter.New(os.Stdout, app)
	logger := log.New(&writer, "Background:  ", log.Default().Flags())

	logger.Print("Hello world!")

	txnName := "logsInContext Sample Transaction"
	txn := app.StartTransaction("logsInContext")

	// Always create a new log object in order to avoid changing the context of the logger for
	// other threads that may be logging outside of this transaction
	txnLogger := log.New(writer.WithTransaction(txn), "Transaction: ", log.Default().Flags())
	txnLogger.Printf("In transaction %s.", txnName)

	// Random sleep to simulate delays
	randomDelay := rand.Intn(300)
	time.Sleep(time.Duration(randomDelay) * time.Millisecond)

	txnLogger.Printf("Ending transaction %s.", txnName)
	txn.End()

	logger.Print("Goodbye!")

	// Wait for shut down to ensure data gets flushed
	app.Shutdown(10 * time.Second)
}
