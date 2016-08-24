package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	flagUpdateURL      string // url to send status to
	flagListenHostPort string // host:port to listen on
)

func init() {
	flag.StringVar(&flagUpdateURL, "update-url", "", "Url to upload status to")
	flag.StringVar(&flagListenHostPort, "listen", "", "host:port to listen on")
}

func main() {
	flag.Parse()

	if flagUpdateURL == "" {
		fmt.Println("Missing required argument: -update-url")
		os.Exit(1)
	}
	if flagListenHostPort == "" {
		fmt.Println("Missing required argument: -listen")
		os.Exit(1)
	}

	// create the DoorServer
	doorServer, err := NewDoorServer(flagUpdateURL)
	if err != nil {
		fmt.Printf("Error creating new DoorServer: %s\n", err)
		os.Exit(1)
	}

	// create our state tracker & change callback
	doorState, err := NewDoorState(func(oldState string, newState string) {
		if oldState != newState {
			fmt.Printf("State changed - oldState: %s; newState: %s\n", oldState, newState)

			err := doorServer.UpdateDoorStatus(newState)
			if err != nil {
				// TODO: keep track of this failure, so we report next time we receive state
				fmt.Printf("Error updating door state with DoorServer: %s\n", err)
			}
		}
	})
	if err != nil {
		fmt.Printf("Error creating DoorState: %s\n", err)
		os.Exit(1)
	}

	// handlers for updating state
	http.HandleFunc("/isopen", func(w http.ResponseWriter, r *http.Request) {
		doorState.RecordState("open")
	})
	http.HandleFunc("/isclosed", func(w http.ResponseWriter, r *http.Request) {
		doorState.RecordState("closed")
	})

	// handler for a crappy HTML interface for a phone
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		header(w)

		// get the state
		state, err := doorState.GetState()
		if err != nil {
			fmt.Fprintf(w, "Error fetching state")
			return
		}

		// output some basic HTML
		updatedAgo := time.Now().Round(time.Second).Sub(state.LastStatusTime.Round(time.Second))
		updatedAgoString := fmt.Sprintf("%s ago", updatedAgo.String())
		if updatedAgo == 0 {
			updatedAgoString = "now"
		}
		if updatedAgo > 5*time.Second {
			fmt.Fprintf(w, "<h1>I DON'T KNOW!</h1>")
			fmt.Fprintf(w, "<h2>Haven't heard from Garage door in %s</h2>", updatedAgo.String())
		} else {
			fmt.Fprintf(w, "<h1>%s</h1>", state.State)
			fmt.Fprintf(w, "<h2>For: %s</h2>", time.Now().Round(time.Second).Sub(state.LastStatusChangeTime.Round(time.Second)).String())
			fmt.Fprintf(w, "<h2>Last updated: %s</h2>", updatedAgoString)
		}
		footer(w)
	})

	err = http.ListenAndServe(flagListenHostPort, nil)
	if err != nil {
		fmt.Printf("Error listening on %s: %s\n", flagListenHostPort, err)
		os.Exit(1)
	}
}

// header writes an HTML header to the writer
func header(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, `
<html><header>
<style type="text/css">
	body {
		font-family: 'Lato', sans-serif;
	}
	h1 {
		font-size: 10rem;
	}
	h2 {
		font-size: 5rem;
	}
</style>
</header>
<body>
`)
}

// footer writes an HTML footer to the writer
func footer(w http.ResponseWriter) {
	fmt.Fprintf(w, "</body></html>")
}
