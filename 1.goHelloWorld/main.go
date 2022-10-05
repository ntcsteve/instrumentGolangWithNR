package main

import (
	"fmt"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	// Making app and err a global variable
	nrApp *newrelic.Application
	nrErr error
)

func main() {
	nrApp, nrErr = newrelic.NewApplication(

		// Workshop > Set the values in the newrelic.Config struct from within a custom newrelic.ConfigOption
		// https://docs.newrelic.com/docs/apm/agents/go-agent/configuration/go-agent-configuration/

		// Workshop > Name your application
		newrelic.ConfigAppName(os.Getenv("APP_NAME")),
		// Workshop > Fill in your New Relic Ingest license key
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		// Workshop > Add debug logging for extra details
		newrelic.ConfigDebugLogger(os.Stdout),
	)

	// If an application could not be created then err will reveal why.
	if nrErr != nil {
		fmt.Println("unable to start NR instrumentation - ", nrErr)
	}

	// Wait for go-agent to avoid data loss
	nrApp.WaitForConnection(5 * time.Second)

	// Run a simple function
	print()

	// Wait for shut down to ensure data gets flushed
	nrApp.Shutdown(5 * time.Second)
}

func print() {
	// Workshop > Monitor a Golang transaction
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-transactions/#go-txn
	nrTxnTracer := nrApp.StartTransaction("Print")
	defer nrTxnTracer.End()

	fmt.Println("Hello world! Welcome to your first instrumented Golang App!")
}
