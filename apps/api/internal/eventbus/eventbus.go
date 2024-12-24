package eventbus

import "sync"

type Event struct {
	Topic string
	Data  interface{}
}

type eventChannel chan Event

type channelSlice []eventChannel

type EventBus struct {
	subscribers map[string]channelSlice
	mu          sync.RWMutex
}

func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string]channelSlice),
	}
}

func (e *EventBus) Subscribe(topic string) <-chan Event {
	ch := make(chan Event)

	e.mu.Lock()
	e.subscribers[topic] = append(e.subscribers[topic], ch)
	e.mu.Unlock()

	return ch
}

func (e *EventBus) Publish(event Event) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, ch := range e.subscribers[event.Topic] {
		go func(c chan Event) {
			c <- event
		}(ch)
	}
}

func (e *EventBus) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, subscribers := range e.subscribers {
		for _, ch := range subscribers {
			close(ch)
		}
	}
}
