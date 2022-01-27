package analytics

import (
	"bytes"
	"sync"

	"github.com/bitrise-io/go-utils/v2/log"
)

const poolSize = 10
const bufferSize = 100

// Tracker ...
type Tracker interface {
	Enqueue(event Event)
	Wait()
}

type tracker struct {
	jobs             chan *bytes.Buffer
	waitGroup        *sync.WaitGroup
	client           Client
	sharedProperties []Property
}

// NewDefaultTracker ...
func NewDefaultTracker(logger log.Logger, shared ...Property) Tracker {
	return NewTracker(NewDefaultClient(logger), shared...)
}

// NewTracker ...
func NewTracker(client Client, shared ...Property) Tracker {
	t := tracker{jobs: make(chan *bytes.Buffer, bufferSize), waitGroup: &sync.WaitGroup{}, client: client, sharedProperties: shared}
	t.init(poolSize)
	return &t
}

// Enqueue ...
func (t tracker) Enqueue(event Event) {
	var b bytes.Buffer
	event.toJson(&b, t.sharedProperties...)
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
