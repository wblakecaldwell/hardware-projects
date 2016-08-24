package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// doorStatusRequest is what we send to AWS
type doorStatusRequest struct {
	Status string `json:"status"` // "open", "closed"
}

// DoorServer is the interface to our AWS endpoints
type DoorServer struct {
	UpdateURL string
}

// NewDoorServer returns a new DoorServer
func NewDoorServer(UpdateURL string) (*DoorServer, error) {
	return &DoorServer{
		UpdateURL: UpdateURL,
	}, nil
}

// UpdateDoorStatus sends the current status of the door to the Door Server
func (ds *DoorServer) UpdateDoorStatus(status string) error {
	// build JSON request body
	jsonBytes, err := json.Marshal(doorStatusRequest{
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("Couldn't marshal DoorStatusRequest: ", err)
	}

	// build HTTP request
	req, err := http.NewRequest("PUT", ds.UpdateURL, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("Couldn't create PUT request for %s: %s", ds.UpdateURL, err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Could't send request to %s: %s", ds.UpdateURL, err)
	}
	defer resp.Body.Close()

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading the response body: %s\n", err)
	}
	fmt.Println("Submitted response to the server. Response body:", string(body))

	return nil
}
