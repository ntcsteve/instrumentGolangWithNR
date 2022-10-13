package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	// Workshop > nrgin integration package
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

// game represents data about a record game.
type game struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	studio string  `json:"studio"`
	Price  float64 `json:"price"`
}

// games slice to seed record game data.
var games = []game{
	{ID: "1", Title: "X-COM", studio: "2K Games", Price: 56.99},
	{ID: "2", Title: "Counter Strike", studio: "Valve", Price: 17.99},
	{ID: "3", Title: "DOOM", studio: "ID Sofrware", Price: 39.99},
	{ID: "4", Title: "Final Fantasy 7", studio: "Square Enix", Price: 69.99},
	{ID: "5", Title: "Command and Conquer", studio: "EA", Price: 19.99},
}

var (
	// Making app and err a global variable
	nrApp *newrelic.Application
	nrErr error
)

func main() {

	nrApp, nrErr = newrelic.NewApplication(
		newrelic.ConfigAppName(os.Getenv("APP_NAME")),
		newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
		// newrelic.ConfigDebugLogger(os.Stdout),
	)

	// If an application could not be created then err will reveal why.
	if nrErr != nil {
		fmt.Println("unable to start NR instrumentation - ", nrErr)
	}

	// Not necessary for monitoring a production application with a lot of data.
	nrApp.WaitForConnection(5 * time.Second)

	router := gin.Default()

	// Workshop > Package nrgin instruments https://github.com/gin-gonic/gin applications.
	// https://pkg.go.dev/github.com/newrelic/go-agent/v3/integrations/nrgin#section-readme
	router.Use(nrgin.Middleware(nrApp))

	router.GET("/games", getgames)
	router.GET("/games/:id", getgameByID)
	router.POST("/games", postgames)

	router.Run("localhost:8080")

	// Wait for shut down to ensure data gets flushed
	nrApp.Shutdown(5 * time.Second)
}

// getgames responds with the list of all games as JSON.
func getgames(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, games)
}

// postgames adds an game from JSON received in the request body.
func postgames(c *gin.Context) {
	var newgame game

	// Call BindJSON to bind the received JSON to newgame.
	if err := c.BindJSON(&newgame); err != nil {
		return
	}

	// Add the new game to the slice.
	games = append(games, newgame)
	c.IndentedJSON(http.StatusCreated, newgame)
}

// getgameByID locates the game whose ID value matches the id
// parameter sent by the client, then returns that game as a response.
func getgameByID(c *gin.Context) {
	id := c.Param("id")

	// Loop through the list of games, looking for an game whose ID value matches the parameter.
	for _, a := range games {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "game not found"})
}
