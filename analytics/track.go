package analytics

import (
	"bytes"

	"github.com/bitrise-io/go-utils/v2/log"
)

const poolSize = 10
const bufferSize = 100

// Tracker ...
type Tracker interface {
	Enqueue(event Event)
	Fork(properties ...Property) Tracker
}

type tracker struct {
	worker     Worker
	Properties []Property
}

// NewDefaultTracker ...
func NewDefaultTracker(logger log.Logger, properties ...Property) Tracker {
	return NewTracker(NewWorker(NewDefaultClient(logger)), properties...)
}

// NewTracker ...
func NewTracker(worker Worker, properties ...Property) Tracker {
	t := tracker{worker: worker, Properties: properties}
	return &t
}

// Enqueue ...
func (t tracker) Enqueue(event Event) {
	var b bytes.Buffer
	event.toJSON(&b, t.Properties...)
	t.worker.Run(&b)
}

// Fork ...
func (t tracker) Fork(properties ...Property) Tracker {
	return NewTracker(t.worker, append(t.Properties, properties...)...)
}
