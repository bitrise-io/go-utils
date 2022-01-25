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

	tracker := NewTrackerBuilder(mockClient).Build()
	tracker.Enqueue(NewEventBuilder("first").Build())
	tracker.Enqueue(NewEventBuilder("second").Build())
	tracker.Enqueue(NewEventBuilder("third").Build())
	tracker.Enqueue(NewEventBuilder("fourth").Build())
	tracker.Enqueue(NewEventBuilder("fifth").Build())
	tracker.Wait()

	mockClient.AssertNumberOfCalls(t, "Send", 5)
}

func Test_tracker_SendIsCalledWithExpectedData(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	tracker := NewTrackerBuilder(mockClient).AddSharedProperty("session", "id").Build()
	tracker.Enqueue(
		NewEventBuilder("first").
			AddProperty("property", "value").
			AddProperty("property2", map[string]string{"foo": "bar"}).
			Build())
	tracker.Wait()

	matcher := mock.MatchedBy(func(buffer *bytes.Buffer) bool {
		var event event
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
