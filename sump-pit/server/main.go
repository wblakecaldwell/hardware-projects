package main

import (
	"fmt"
	"github.com/wblakecaldwell/profiler"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	global_error error
)

func main() {
	database := NewMemoryDatabase()

	// load config
	config, err := NewConfigFromEnv()
	if err != nil {
		fmt.Printf("Error parsing configuration from global variables: %s\n", err)
		global_error = err
		return
	}

	// emailer is responsible for keeping track of how often to email the report, how to send it, and whom to send it to
	emailer := NewSMTPEmailer(
		config.EmailSMTPLogin,
		config.EmailSMTPPassword,
		config.EmailSMTPServer,
		config.EmailSMTPPort,
		config.EmailSender,
		config.EmailRecipients,
		1*time.Minute)

	// all locking done at the handler level
	rwMutex := sync.RWMutex{}

	// standard endpoints
	http.HandleFunc("/", indexHtmlHandler)
	http.HandleFunc("/index.html", indexHtmlHandler)
	http.HandleFunc("/info", buildSumpInfoHandler(database, 2*time.Hour, &rwMutex))
	http.HandleFunc("/water-level", buildSumpRegisterLevelsHandler(database, config.PanicWaterLevel, emailer, config.ServerSecret, &rwMutex))

	// profiler endpoints
	profiler.AddMemoryProfilingHandlers()

	err = http.ListenAndServe(config.ListenHostPort, nil)
	if err != nil {
		fmt.Printf("Error listening on %s: %s", config.ListenHostPort, err)
		os.Exit(1)
	}
}
