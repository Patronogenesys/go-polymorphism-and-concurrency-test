package clock

import (
	"modellingSystems/devicefailureExperiment/config"
	"modellingSystems/devicefailureExperiment/models"
)

type DefaultExperimentClockFactory struct {
}

func (f DefaultExperimentClockFactory) New() *ExperimentClock {
	return &ExperimentClock{
		tickChan:   make(chan models.ExperimentTickEvent),
		stopChan:   make(chan struct{}, 1),
		Time:       0,
		StopTime:   config.DefaultExperimentDuration,
		WasStarted: false,
	}
}

type ExperimentClockFactory struct {
	Duration int
}

func (f ExperimentClockFactory) New() *ExperimentClock {
	return &ExperimentClock{
		tickChan:   make(chan models.ExperimentTickEvent),
		stopChan:   make(chan struct{}, 1),
		Time:       0,
		StopTime:   f.Duration,
		WasStarted: false,
	}
}
