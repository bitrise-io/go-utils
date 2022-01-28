package analytics

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/bitrise-io/go-utils/v2/analytics/mocks"
	"github.com/stretchr/testify/mock"
)

func Test_tracker_EnqueueWaitCycleExecutesSends(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	worker := NewWorker(mockClient)
	tracker := NewTracker(worker)
	tracker.Enqueue(NewEvent("first"))
	tracker.Enqueue(NewEvent("second"))
	tracker.Enqueue(NewEvent("third"))
	tracker.Enqueue(NewEvent("fourth"))
	tracker.Enqueue(NewEvent("fifth"))
	worker.Wait()

	mockClient.AssertNumberOfCalls(t, "Send", 5)
}

func Test_tracker_SendIsCalledWithExpectedData(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	worker := NewWorker(mockClient)
	tracker := NewTracker(worker, StringProperty("session", "id"))
	tracker.Enqueue(NewEvent(
		"first",
		StringProperty("property", "value"),
		IntProperty("intproperty", 42),
		LongProperty("longproperty", 42),
		FloatProperty("floatproperty", 3.14),
		BoolProperty("boolproperty", true),
		NestedProperty("property2", StringProperty("foo", "bar"))),
	)
	worker.Wait()

	matcher := mock.MatchedBy(func(buffer *bytes.Buffer) bool {
		var event event
		err := json.Unmarshal(buffer.Bytes(), &event)
		if err != nil {
			return false
		}
		if event.EventName != "first" {
			return false
		}
		if len(event.Properties) != 7 ||
			event.Properties["property"] != "value" ||
			event.Properties["intproperty"].(float64) != 42 ||
			event.Properties["longproperty"].(float64) != 42 ||
			event.Properties["floatproperty"].(float64) != 3.14 ||
			event.Properties["boolproperty"].(bool) != true ||
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

func Test_tracker_ForkedTrackerAddsProperties(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	worker := NewWorker(mockClient)
	tracker := NewTracker(worker, StringProperty("first", "first"))
	tracker.Enqueue(NewEvent("event"))
	forked := tracker.Fork(StringProperty("second", "second"))
	forked.Enqueue(NewEvent("event2"))
	worker.Wait()

	mockClient.AssertNumberOfCalls(t, "Send", 2)
	matcher := mock.MatchedBy(func(buffer *bytes.Buffer) bool {
		var event event
		err := json.Unmarshal(buffer.Bytes(), &event)
		if err != nil {
			return false
		}
		if event.EventName == "event" {
			if len(event.Properties) != 1 || event.Properties["first"] != "first" {
				return false
			}
			return true
		}
		if event.EventName == "event2" {
			if len(event.Properties) != 2 || event.Properties["first"] != "first" || event.Properties["second"] != "second" {
				return false
			}
			return true
		}
		return false
	})
	mockClient.AssertCalled(t, "Send", matcher)
}
