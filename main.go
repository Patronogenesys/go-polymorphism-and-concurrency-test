package main

import (
	"fmt"
	"modellingSystems/devicefailureExperiment/models/experiment/strategies"
	"modellingSystems/devicefailureExperiment/models/simulation"
)

// There i tried to plan architecture

// Domain system constants
// Replacement cost						INT  	-> Device
// Await cost, Await time 				INT 	-> Experiment
// Probability of failure over time 	FUNC 	-> Device
// Replacement strategy					FUNC 	-> Experiment
// Device group size					INT 	-> Experiment	Represented by len(Experiment.Devices)

// Representative system State
// Device working time					INT	-> Device
// Time									INT	-> Experiment
// Last failure time					INT	-> Experiment

// Stats
// Money spent on replacement					-> Experiment
// Money spent on await							-> Experiment
// Number of failed devices						-> Experiment

// Rules
// Device fails at time t if predicate(time) is true
// After device fails, it is replaced over AwaitTime, and AwaitCost increments by ReplacementCost
// While device is awaiting, it is not available for failure

// Experiment is a state machine
// States
// 1. Uninitialized
// 2. Working											Represented
// 3. Replacing
// 4. Exited
// Transitions
// 1. Uninitialized -> (onInit) -> Working 				Represented by method Factory.Device() and go Experiment.Rum()
// 2. Working -> (onDeviceFail) -> Replacing			Represented by event DeviceFail event
// 3. Replacing -> (onReplaced) -> Working				Represented by event DeviceReplaced
// 2. Working -> (onExit) -> Exited						Represented by event Tick (type: Stop)

// Events in system (Event: Sender -> Receivers)
// Tick: Simulation -> Experiment -> Device	 		  	Represented by chan Experiment.TickChan, chan Device.TickChan,
//														type ExperimentTickEvent
//															(type ExperimentTickEventType: <TimeIncrement | Stop>)),
//														type DeviceTickEvent
//
// DeviceFail: Device -> Experiment						Represented by chan Device.FailChan := &Experiment.FailChan
// 														type DeviceFailEvent
// DeviceReplaced: Experiment -> Experiment				Represented by chan Experiment.ReplacedChan
// ExperimentExited: Experiment -> Simulation			Represented by chan Experiment.ResultChan

func main() {
	// Set parameters at [./config/defaults.go]

	// test the Experiment
	resultChan := make(chan simulation.Result)
	simFactory := &simulation.DefaultFactory{
		Name:                "ReplaceOnlyFailed",
		ReplacementStrategy: strategies.ReplaceOnlyFailed,
		ResultChan:          resultChan,
	}
	sim := simFactory.New()
	sim.Run()
	fmt.Println(<-resultChan)

	simFactory = &simulation.DefaultFactory{
		Name:                "ReplaceFailedAndTheOldest",
		ReplacementStrategy: strategies.ReplaceFailedAndTheOldest,
		ResultChan:          resultChan,
	}
	sim = simFactory.New()
	sim.Run()
	fmt.Println(<-resultChan)
}
