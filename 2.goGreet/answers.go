package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

var (
	nrApp *newrelic.Application
	nrErr error
)

func main() {

	nrApp, nrErr = newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		newrelic.ConfigDebugLogger(os.Stdout),

		// Workshop > Additional configuration via a config function
		// https://docs.newrelic.com/docs/apm/agents/go-agent/configuration/go-agent-configuration
		func(config *newrelic.Config) {
			config.Enabled = true
			config.DistributedTracer.Enabled = true
			config.Labels = map[string]string{
				"Env":      "Dev",
				"Function": "Greetings",
				"Platform": "My Machine",
				"Team":     "API",
			}
		},
	)

	// If an application could not be created then err will reveal why.
	if nrErr != nil {
		fmt.Println("unable to start NR instrumentation - ", nrErr)
	}

	// Not necessary for monitoring a production application with a lot of data.
	nrApp.WaitForConnection(5 * time.Second)

	// Request a greeting message.
	message, err := Hello("Tony Stark")

	// If an error was returned, print it to the console and exit the program.
	if err != nil {
		log.Fatal(err)
	}

	// If no error was returned, print the returned message to the console.
	fmt.Println(message)

	// Wait for shut down to ensure data gets flushed
	nrApp.Shutdown(5 * time.Second)
}

// init sets initial values for variables used in the function.
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Hello returns a greeting for the named person.
func Hello(name string) (string, error) {

	// Monitor a transaction
	nrTxnTracer := nrApp.StartTransaction("Hello")
	defer nrTxnTracer.End()

	// If no name was given, return an error with a message.
	if name == "" {
		return name, errors.New("empty name")
	}

	// Create a message using a random format.
	message := fmt.Sprintf(randomFormat(), name)

	// Workshop > Custom attributes by using this method in a transaction
	// https://docs.newrelic.com/docs/apm/agents/go-agent/api-guides/guide-using-go-agent-api#metadata
	nrTxnTracer.AddAttribute("message", message)
	return message, nil
}

// randomFormat returns one of a set of greeting messages.
// The returned message is selected at random.
func randomFormat() string {

	// Workshop > Monitor a transaction
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-transactions/#go-txn
	nrTxnTracer := nrApp.StartTransaction("randomFormat")
	defer nrTxnTracer.End()

	// Random sleep to simulate delays
	randomDelayOuter := rand.Intn(40)
	time.Sleep(time.Duration(randomDelayOuter) * time.Microsecond)

	// Workshop > Create a segment
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-segments
	nrSegment := nrTxnTracer.StartSegment("Formats")

	// Random sleep to simulate delays
	randomDelayInner := rand.Intn(80)
	time.Sleep(time.Duration(randomDelayInner) * time.Microsecond)

	// A slice of message formats.
	formats := []string{
		"Hi, %v. Welcome!",
		"Great to see you, %v!",
		"Good day, %v! Well met!",
		"%v! Hi there!",
		"Greetings %v!",
		"Hello there, %v!",
	}

	// Workshop > End a segment
	// https://docs.newrelic.com/docs/apm/agents/go-agent/instrumentation/instrument-go-segments
	nrSegment.End()

	// Return a randomly selected message format by specifying
	// a random index for the slice of formats.
	return formats[rand.Intn(len(formats))]
}
