package clock

import (
	"errors"
	"modellingSystems/devicefailureExperiment/models"
)

type Clock interface {
	Start() error
	Stop() // Note: Possible bug multiple stops scenario
	TickChan() <-chan models.ExperimentTickEvent
	tick()
}

type ExperimentClock struct {
	tickChan chan models.ExperimentTickEvent
	// Use this channel to stop the clock
	stopChan   chan struct{}
	Time       int
	StopTime   int
	WasStarted bool
}

// In every tick, the clock increments time by 1

func (c *ExperimentClock) tick() {
	c.Time++
	if c.Time >= c.StopTime {
		c.Stop()
		return
	}
	c.tickChan <- models.ExperimentTickEvent{Type: models.TimeIncrement}
}

// Start the clock goroutine
func (c *ExperimentClock) Start() error {
	if c.WasStarted {
		return errors.New("clock already started")
	}
	go func() {
		for {
			select {
			case <-c.stopChan:
				return
			default:
				c.tick()
			}
		}
	}()

	c.WasStarted = true
	return nil
}

func (c *ExperimentClock) Stop() {
	// stop ticking goroutine
	c.stopChan <- struct{}{}
	// send a stop signal to the tick channel
	c.tickChan <- models.ExperimentTickEvent{Type: models.Stop}
	close(c.tickChan)
}

func (c *ExperimentClock) TickChan() <-chan models.ExperimentTickEvent {
	return c.tickChan
}
