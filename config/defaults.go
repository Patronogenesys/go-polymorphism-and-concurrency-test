package config

import (
	"math"
)

// Device
var (
	ReplacementCostDefault     = 50
	ReplacementCostMaskDefault = []int{50, 100, 300}

	DefaultLambda               = 0.1
	ProbabilityOfFailureDefault = func(time int) float64 {
		return 1 - math.Exp(-DefaultLambda*float64(time))
	}
)

// Experiment

var (
	DeviceNumberDefault    = 3
	AwaitPriceDefault      = 40
	ReplaceDurationDefault = 2
)

// ExperimentClock

var DefaultExperimentDuration = 100

// Simulation

var ExperimentNumberDefault int = 1e6
