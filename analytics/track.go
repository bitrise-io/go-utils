package analytics

import (
	"bytes"
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

const poolSize = 10
const bufferSize = 100

// Tracker ...
type Tracker interface {
	Enqueue(eventName string, properties map[string]interface{})
	Wait()
}

type tracker struct {
	jobs             chan *bytes.Buffer
	waitGroup        *sync.WaitGroup
	client           Client
	sharedProperties map[string]interface{}
}

// NewDefaultTracker ...
func NewDefaultTracker(shared map[string]interface{}) Tracker {
	return NewTracker(make(chan *bytes.Buffer, bufferSize), &sync.WaitGroup{}, poolSize, NewDefaultClient(), shared)
}

// NewTracker ...
func NewTracker(jobs chan *bytes.Buffer, waitGroup *sync.WaitGroup, size int, client Client, sharedProperties map[string]interface{}) Tracker {
	t := tracker{jobs: jobs, waitGroup: waitGroup, client: client, sharedProperties: sharedProperties}
	t.init(size)
	return &t
}

// Enqueue ...
func (t tracker) Enqueue(eventName string, properties map[string]interface{}) {
	mergedProperties := make(map[string]interface{})
	for k, v := range t.sharedProperties {
		mergedProperties[k] = v
	}
	for k, v := range properties {
		mergedProperties[k] = v
	}
	event := Event{
		ID:         uuid.NewString(),
		EventName:  eventName,
		Timestamp:  time.Now().UnixNano(),
		Properties: mergedProperties,
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(event); err != nil {
		panic("Analytics event should be serializable to JSON")
	}
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
