package analytics

import (
	"bytes"
	"encoding/json"
	"github.com/bitrise-io/go-utils/v2/analytics/mocks"
	"github.com/stretchr/testify/mock"
	"sync"
	"testing"
)

func Test_tracker_EnqueueWaitCycleExecutesSends(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	tracker := NewTracker(make(chan *bytes.Buffer, bufferSize), &sync.WaitGroup{}, 5, mockClient, map[string]interface{}{})
	tracker.Enqueue("first", map[string]interface{}{})
	tracker.Enqueue("second", map[string]interface{}{})
	tracker.Enqueue("third", map[string]interface{}{})
	tracker.Enqueue("fourth", map[string]interface{}{})
	tracker.Enqueue("fifth", map[string]interface{}{})
	tracker.Wait()

	mockClient.AssertNumberOfCalls(t, "Send", 5)
}

func Test_tracker_SendIsCalledWithExpectedData(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	tracker := NewTracker(make(chan *bytes.Buffer, bufferSize), &sync.WaitGroup{}, 5, mockClient, map[string]interface{}{"session": "id"})
	tracker.Enqueue("first", map[string]interface{}{
		"property":  "value",
		"property2": map[string]string{"foo": "bar"},
	})
	tracker.Wait()

	matcher := mock.MatchedBy(func(buffer *bytes.Buffer) bool {
		var event Event
		err := json.Unmarshal(buffer.Bytes(), &event)
		if err != nil {
			return false
		}
		if event.EventName != "first" {
			return false
		}
		if len(event.Properties) != 3 ||
			event.Properties["property"] != "value" ||
			event.Properties["session"] != "id" ||
			event.Properties["property2"].(map[string]interface{})["foo"] != "bar" {
			return false
		}
		if event.ID == "" || event.Timestamp == 0 {
			return false
		}
		return true
	})
	mockClient.AssertNumberOfCalls(t, "Send", 1)
	mockClient.AssertCalled(t, "Send", matcher)
}
