package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
)

// init sets initial values for variables used in the function.
func init() {
	rand.Seed(time.Now().UnixNano())
}

func async(w http.ResponseWriter, r *http.Request) {
	// To access the transaction in your handler, use the newrelic.FromContext API.
	txn := newrelic.FromContext(r.Context())
	// This WaitGroup is used to wait for all the goroutines to finish.
	wg := &sync.WaitGroup{}
	println("goRoutines created!")

	for i := 1; i < 9; i++ {
		wg.Add(1)
		i := i

		// Workshop > trace asynchronous applications
		// The Transaction.NewGoroutine() allows transactions to create segments in multiple goroutines.
		// https://docs.newrelic.com/docs/apm/agents/go-agent/features/trace-asynchronous-applications
		go func(txn *newrelic.Transaction) {
			defer wg.Done()
			defer txn.StartSegment("goroutine" + strconv.Itoa(i)).End()
			println("goRoutine " + strconv.Itoa(i))

			randomDelay := rand.Intn(500)
			time.Sleep(time.Duration(randomDelay) * time.Millisecond)
		}(txn.NewGoroutine())
	}

	// Workshop > Ensure the WaitGroup is done
	segment := txn.StartSegment("WaitGroup")
	wg.Wait()
	segment.End()
	w.Write([]byte("success!"))
}

func main() {
	nrApp, nrErr := newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		// newrelic.ConfigDebugLogger(os.Stdout),
	)
	if nrErr != nil {
		fmt.Println(nrErr)
		os.Exit(1)
	}

	// Wait for shut down to ensure data gets flushed
	nrApp.WaitForConnection(5 * time.Second)

	// Workshop > ListenAndServe starts an HTTP server with a given address and handler
	// http.HandleFunc(newrelic.WrapHandleFunc(nrApp, "/async", async))
	// http.ListenAndServe(":8000", nil)

	// Wait for shut down to ensure data gets flushed
	nrApp.Shutdown(5 * time.Second)
}
