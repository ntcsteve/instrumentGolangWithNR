package main

import (
	"fmt"
	"io"
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

// Workshop > You may track errors using the Transaction.NoticeError method.
// The easiest way to get started with NoticeError is to use errors based on Go's standard error interface.
// https://github.com/newrelic/go-agent/blob/master/GUIDE.md#error-reporting
func noticeErrorWithAttributes(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "noticing an error")

	txn := newrelic.FromContext(r.Context())
	txn.NoticeError(newrelic.Error{
		Message: "uh oh. something went very wrong",
		Class:   "errors are aggregated by class",
		Attributes: map[string]interface{}{
			"important_number": 97232,
			"relevant_string":  "classError",
		},
	})
	println("Oops, there is an error!")
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
	w.Write([]byte("goRoutines success!"))
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
	// http.HandleFunc(newrelic.WrapHandleFunc(nrApp, "/error", noticeErrorWithAttributes))
	// http.HandleFunc(newrelic.WrapHandleFunc(nrApp, "/async", async))
	http.ListenAndServe(":8000", nil)

	// Wait for shut down to ensure data gets flushed
	nrApp.Shutdown(5 * time.Second)
}
