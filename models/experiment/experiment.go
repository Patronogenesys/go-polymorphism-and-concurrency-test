package experiment

import (
	"errors"
	"fmt"
	"modellingSystems/devicefailureExperiment/models"
	"modellingSystems/devicefailureExperiment/models/clock"
	"modellingSystems/devicefailureExperiment/models/device"
	"modellingSystems/devicefailureExperiment/models/experiment/strategies"
	"sync"
)

type Interface interface {
	Run(group *sync.WaitGroup) <-chan Result
}

type State int

const (
	Working State = iota
	Replacing
)

type Result struct {
	Error            error
	ReplacingExpense int
	AwaitExpense     int
	FailedDevices    int
}

func (r Result) String() string {
	return fmt.Sprintf("Error: %v, ReplacingExpense: %v, AwaitExpense: %v, FailedDevices: %v",
		r.Error,
		r.ReplacingExpense,
		r.AwaitExpense,
		r.FailedDevices,
	)
}

type Experiment struct {
	Devices            []device.Device
	DeviceArrayFactory models.ArrayFactory[device.Device]
	Clock              clock.Clock

	// State
	CurrentStateDuration int
	state                State

	// Constants (Inputs)
	AwaitPrice          int
	ReplaceDuration     int
	ReplacementStrategy strategies.Strategy

	// Stats
	ReplacingExpense int
	AwaitExpense     int
	FailedDevices    int
	ReplacedDevices  int

	// Output
	ResultChan chan Result
}

func (e *Experiment) Run(group *sync.WaitGroup) <-chan Result {
	go func() {
		defer group.Done()
		// Start the Clock
		if e.Clock.Start() != nil {
			e.ResultChan <- Result{Error: errors.New("clock already started")}
			return
		}
		// While Clock is running, handle ticks
		for tick := range e.Clock.TickChan() {
			e.handleTick(tick)
		}
	}()
	return e.ResultChan
}

func (e *Experiment) handleTick(tick models.ExperimentTickEvent) {
	// Handle Stop tick
	if tick.Type == models.Stop {
		e.ResultChan <- Result{Error: nil, ReplacingExpense: e.ReplacingExpense,
			AwaitExpense: e.AwaitExpense, FailedDevices: e.FailedDevices}
	}
	// Handle TimeIncrement tick
	if tick.Type == models.TimeIncrement {
		// State machine
		if e.state == Working {
			e.handleWorkingState()
		}
		if e.state == Replacing {
			e.handleReplacingState()
		}
	}
}

// Working state IncrementTick handler

func (e *Experiment) handleWorkingState() {
	e.CurrentStateDuration++
	for i := range e.Devices {
		e.Devices[i].WorkingTime++
		if e.Devices[i].Failed() {
			e.onDeviceFail()
		}
	}
}

// Working -> Replacing Transition

func (e *Experiment) onDeviceFail() {
	e.state = Replacing
	e.CurrentStateDuration = 0
}

// Replacing state IncrementTick handler

func (e *Experiment) handleReplacingState() {
	e.CurrentStateDuration++
	if e.CurrentStateDuration == e.ReplaceDuration {
		e.onReplaced()
	}
}

// Replacing -> Working Transition

func (e *Experiment) onReplaced() {
	// Count failed devices
	for i := range e.Devices {
		if e.Devices[i].WasFailed {
			e.FailedDevices++
		}
	}
	// Replace devices

	wasReplaced := e.ReplacementStrategy(e.Devices, e.DeviceArrayFactory)
	for i := range wasReplaced {
		if wasReplaced[i] {
			e.applyReplacePenalties(&e.Devices[i])
		}
	}
	e.applyAwaitPenalties()
	// Reset state
	e.state = Working
	e.CurrentStateDuration = 0
}

// Each replaced device costs ReplacePrice

func (e *Experiment) applyReplacePenalties(d *device.Device) {
	e.ReplacingExpense += d.ReplacementCost
	e.ReplacedDevices++
}

// Independent of the number of replaced devices, each onReplaced transition costs AwaitPrice

func (e *Experiment) applyAwaitPenalties() {
	e.AwaitExpense += e.AwaitPrice
}
