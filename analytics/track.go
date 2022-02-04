package analytics

import (
	"bytes"
	"github.com/bitrise-io/go-utils/v2/log"
	"sync"
)

const poolSize = 10
const bufferSize = 100

// Properties ...
type Properties map[string]interface{}

// Merge ...
func (p Properties) Merge(properties Properties) Properties {
	r := Properties{}
	for key, value := range p {
		r[key] = value
	}
	for key, value := range properties {
		r[key] = value
	}
	return r
}

// Tracker ...
type Tracker interface {
	Enqueue(eventName string, properties ...Properties)
	Wait()
}

type tracker struct {
	jobs       chan *bytes.Buffer
	waitGroup  *sync.WaitGroup
	client     Client
	properties []Properties
}

// NewDefaultTracker ...
func NewDefaultTracker(properties ...Properties) Tracker {
	return NewTracker(NewDefaultClient(log.NewLogger()), properties...)
}

// NewTracker ...
func NewTracker(client Client, properties ...Properties) Tracker {
	t := tracker{client: client, jobs: make(chan *bytes.Buffer, bufferSize), waitGroup: &sync.WaitGroup{}, properties: properties}
	t.init(poolSize)
	return &t
}

// Enqueue ...
func (t tracker) Enqueue(eventName string, properties ...Properties) {
	var b bytes.Buffer
	newEvent(eventName, append(t.properties, properties...)).toJSON(&b)
	t.jobs <- &b
	t.waitGroup.Add(1)
}

// Wait ...
func (t tracker) Wait() {
	close(t.jobs)
	t.waitGroup.Wait()
}

func (t tracker) init(size int) {
	for i := 0; i < size; i++ {
		go t.worker()
	}
}

func (t tracker) worker() {
	for job := range t.jobs {
		t.client.Send(job)
		t.waitGroup.Done()
	}
}
