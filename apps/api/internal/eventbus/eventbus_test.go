package eventbus_test

import (
	"rango/api/internal/eventbus"
	"reflect"
	"testing"
)

func TestEventBus(t *testing.T) {
	eventBus := eventbus.New()
	testTopic := "topic.test"
	testData := "test data"

	wantedEvent := eventbus.Event{
		Topic: testTopic,
		Data:  testData,
	}

	subscriber := eventBus.Subscribe(testTopic)

	eventBus.Publish(wantedEvent)

	if actualEvent := <-subscriber; !reflect.DeepEqual(wantedEvent, actualEvent) {
		t.Errorf("Did not get the published event. Wanted: %v. Got: %v", wantedEvent, actualEvent)
	}
}
