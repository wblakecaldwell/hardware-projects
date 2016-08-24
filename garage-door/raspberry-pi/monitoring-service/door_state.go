package main

import (
	"sync"
	"time"
)

// StateChangedFunc defines a method signature for a door state change handler
type StateChangedFunc func(oldState string, newState string)

// DoorStateResponse is used to respond to status queries
type DoorStateResponse struct {
	State                string
	LastStatusTime       time.Time
	LastStatusChangeTime time.Time
}

// DoorState keeps track of the state of the door
type DoorState struct {
	lock                 sync.Mutex
	lastStatusTime       time.Time
	lastStatusChangeTime time.Time
	state                string
	stateChangedFunc     StateChangedFunc
}

// NewDoorState returns a new DoorState
func NewDoorState(stateChangedFunc StateChangedFunc) (*DoorState, error) {
	ds := DoorState{
		lock:             sync.Mutex{},
		stateChangedFunc: stateChangedFunc,
	}
	return &ds, nil
}

// RecordState records the current state, calling the callback function
// if it's changed.
func (ds *DoorState) RecordState(state string) error {
	ds.lock.Lock()

	oldState := ds.state
	ds.state = state
	stateChanged := oldState != state

	ds.lastStatusTime = time.Now()
	if stateChanged {
		ds.lastStatusChangeTime = time.Now()
	}
	ds.lock.Unlock() // need to call the stateChangedFunc outside the lock

	if stateChanged {
		ds.stateChangedFunc(oldState, state)
	}
	return nil
}

// GetState returns a snapshot of the current state
func (ds *DoorState) GetState() (*DoorStateResponse, error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	response := DoorStateResponse{
		State:                ds.state,
		LastStatusTime:       ds.lastStatusTime,
		LastStatusChangeTime: ds.lastStatusChangeTime,
	}
	return &response, nil
}
