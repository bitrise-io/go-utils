package analytics

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/bitrise-io/go-utils/v2/mocks"

	"github.com/stretchr/testify/mock"
)

func Test_tracker_EnqueueWaitCycleExecutesSends(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	tracker := NewTracker(mockClient, timeout)
	tracker.Enqueue("first")
	tracker.Enqueue("second")
	tracker.Enqueue("third")
	tracker.Enqueue("fourth")
	tracker.Enqueue("fifth")
	tracker.Wait()

	mockClient.AssertNumberOfCalls(t, "Send", 5)
}

func Test_tracker_SendIsCalledWithExpectedData(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	tracker := NewTracker(mockClient, timeout)
	baseProperties := Properties{"session": "id"}
	tracker.Enqueue(
		"first",
		baseProperties, Properties{
			"property":      "value",
			"intproperty":   42,
			"longproperty":  42,
			"floatproperty": 3.14,
			"boolproperty":  true,
			"property2":     Properties{"foo": "bar"},
		},
	)
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

func Test_tracker_MergingPropertiesWork(t *testing.T) {
	mockClient := new(mocks.Client)
	mockClient.On("Send", mock.Anything).Return()

	tracker := NewTracker(mockClient, timeout, Properties{"base": "base"})
	baseProperties := Properties{"first": "first"}
	tracker.Enqueue("event", baseProperties)
	newBaseProperties := baseProperties.Merge(Properties{"second": "second"})
	tracker.Enqueue("event2", newBaseProperties)
	tracker.Wait()

	mockClient.AssertNumberOfCalls(t, "Send", 2)
	matcher := mock.MatchedBy(func(buffer *bytes.Buffer) bool {
		var event event
		err := json.Unmarshal(buffer.Bytes(), &event)
		if err != nil {
			return false
		}
		if event.EventName == "event" {
			if len(event.Properties) != 2 || event.Properties["base"] != "base" || event.Properties["first"] != "first" {
				return false
			}
			return true
		}
		if event.EventName == "event2" {
			if len(event.Properties) != 3 || event.Properties["base"] != "base" || event.Properties["first"] != "first" || event.Properties["second"] != "second" {
				return false
			}
			return true
		}
		return false
	})
	mockClient.AssertCalled(t, "Send", matcher)
}

func Test_tracker_WaitTimesOutOnBlockingClient(t *testing.T) {
	timeout := time.After(3 * time.Second)
	done := make(chan bool)
	go func() {
		mockClient := new(mocks.Client)
		mockClient.On("Send", mock.Anything).Run(func(args mock.Arguments) {
			select {}
		})
		tracker := NewTracker(mockClient, 2*time.Second)
		tracker.Enqueue("block")
		tracker.Wait()
		done <- true
	}()

	select {
	case <-timeout:
		t.Fatal("Test didn't finish in time")
	case <-done:
	}
}

func Test_NewDefaultTracker_Disabled(t *testing.T) {
	t.Setenv(analyticsDisabledEnv, "true")

	tracker := NewDefaultTracker(new(mocks.Logger))

	if _, ok := tracker.(noopTracker); !ok {
		t.Fatalf("expected noopTracker when %s is set", analyticsDisabledEnv)
	}
}
