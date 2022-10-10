// Copyright 2020 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

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
	txn := newrelic.FromContext(r.Context())
	wg := &sync.WaitGroup{}

	for i := 1; i < 9; i++ {
		wg.Add(1)
		i := i
		println(i)
		go func(txn *newrelic.Transaction) {
			defer wg.Done()
			defer txn.StartSegment("goroutine" + strconv.Itoa(i)).End()

			randomDelayInner := rand.Intn(800)
			time.Sleep(time.Duration(randomDelayInner) * time.Millisecond)
		}(txn.NewGoroutine())
	}

	segment := txn.StartSegment("WaitGroup")
	wg.Wait()
	segment.End()
	w.Write([]byte("success!"))
}

func main() {
	app, err := newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		// newrelic.ConfigDebugLogger(os.Stdout),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Wait for shut down to ensure data gets flushed
	app.WaitForConnection(5 * time.Second)

	http.HandleFunc(newrelic.WrapHandleFunc(app, "/async", async))
	http.ListenAndServe(":8000", nil)

	// Wait for shut down to ensure data gets flushed
	app.Shutdown(5 * time.Second)
}
