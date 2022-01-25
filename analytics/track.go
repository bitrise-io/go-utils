package analytics

import (
	"bytes"
	"sync"

	"github.com/bitrise-io/go-utils/v2/log"
)

const poolSize = 10
const bufferSize = 100

// TrackerBuilder ...
type TrackerBuilder interface {
	AddSharedProperty(key string, value interface{}) TrackerBuilder
	Build() Tracker
}

type trackerBuilder struct {
	client           Client
	sharedProperties map[string]interface{}
}

// NewTrackerBuilder ...
func NewTrackerBuilder(client Client) TrackerBuilder {
	return trackerBuilder{client: client, sharedProperties: map[string]interface{}{}}
}

// AddSharedProperty ...
func (t trackerBuilder) AddSharedProperty(key string, value interface{}) TrackerBuilder {
	t.sharedProperties[key] = value
	return t
}

// Build ...
func (t trackerBuilder) Build() Tracker {
	return newTracker(t.client, t.sharedProperties)
}

// Tracker ...
type Tracker interface {
	Enqueue(event Event)
	Wait()
}

type tracker struct {
	jobs             chan *bytes.Buffer
	waitGroup        *sync.WaitGroup
	client           Client
	sharedProperties map[string]interface{}
}

// NewDefaultTrackerBuilder ...
func NewDefaultTrackerBuilder(logger log.Logger) TrackerBuilder {
	return NewTrackerBuilder(NewDefaultClient(logger))
}

// newTracker ...
func newTracker(client Client, shared map[string]interface{}) Tracker {
	t := tracker{jobs: make(chan *bytes.Buffer, bufferSize), waitGroup: &sync.WaitGroup{}, client: client, sharedProperties: shared}
	t.init(poolSize)
	return &t
}

// Enqueue ...
func (t tracker) Enqueue(event Event) {
	var b bytes.Buffer
	event.toJson(&b, t.sharedProperties)
	t.jobs <- &b
	t.waitGroup.Add(1)
}

// Wait ...
func (t tracker) Wait() {
	close(t.jobs)
	t.waitGroup.Wait()
}

func (t tracker) init(size int) {
	for w := 0; w < size; w++ {
		go t.worker()
	}
}

func (t tracker) worker() {
	for job := range t.jobs {
		t.client.Send(job)
		t.waitGroup.Done()
	}
}
