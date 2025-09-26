package main

import (
	"fmt"
	"slices"
	"sync"
)

type UnSub func()
type SubCallback[T any] func(T)

type Topic[T any] struct {
	channels []chan T
	mu   sync.Mutex
}

func NewTopic[T any]() *Topic[T] {
	topic := Topic[T]{
		channels: make([]chan T, 0),
		mu:   sync.Mutex{},
	}

	return &topic
}

func (t *Topic[T]) Close() {
	// TODO: cleanup
}

func (t *Topic[T]) Pub(v T) {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, c := range t.channels {
		c <- v
	}
}

func (t *Topic[T]) Sub(cb SubCallback[T]) UnSub {
	sub, unsub := t.SubChan()
	go func() {
		for v := range sub {
			cb(v)
		}
	}()
	return unsub
}

func (t *Topic[T]) SubChan() (chan T, UnSub) {
	c := make(chan T)
	t.channels = append(t.channels , c)

	unsub := func() {
		t.channels = slices.DeleteFunc(t.channels, func(cc chan T) bool {
			return cc == c
		})
	}
	return c, unsub
}

type EventBus struct {
	topics map[string]any
}

func NewEventBus() *EventBus {
	return &EventBus { topics: make(map[string]any) }
}

func CreateTopic[T any](bus *EventBus, topicName string) (*Topic[T], error) {
	_, found := bus.topics[topicName]
	if found {
		return nil, fmt.Errorf("topic %s already exists", topicName)
	} else {
		topic := &Topic[T]{}
		bus.topics[topicName] = topic
		return topic, nil
	}
}

// Get topic from event bus of name `topicName`
// if theres no topic named `topicName` then creates a new topic
func GetTopic[T any](bus *EventBus, topicName string) (*Topic[T], error) {
	var topic *Topic[T]
	topicAny, found := bus.topics[topicName]
	if found {
		if typedTopic, ok := topicAny.(*Topic[T]); ok {
			topic = typedTopic
		} else {
			return nil, fmt.Errorf("topic %s is not type %T", topicName, *new(T))
		}
	} else {
		return nil, fmt.Errorf("topic %s does not exist", topicName)
	}
	return topic, nil
}

func CloseTopic[T any](bus *EventBus, topicName string) error {
	topic, err := GetTopic[T](bus, topicName)
	if err != nil {
		return err
	}
	topic.Close()
	return nil
}

func Pub[T any](bus *EventBus, topicName string, v T) error {
	topic, err := GetTopic[T](bus, topicName)
	if err != nil {
		return err
	}
	topic.Pub(v)
	return nil
}

func Sub[T any](bus *EventBus, topicName string, cb SubCallback[T]) (UnSub, error) {
	topic, err := GetTopic[T](bus, topicName)
	if err != nil {
		return nil, err
	}
	unsub := topic.Sub(cb)
	return unsub, nil
}

func SubChan[T any](bus *EventBus, topicName string) (chan T, UnSub, error) {
	topic, err := GetTopic[T](bus, topicName)
	if err != nil {
		return nil, nil, err
	}
	c, unsub := topic.SubChan()
	return c, unsub, nil
}

var lock sync.Mutex
var bus *EventBus

func GetEventBus() *EventBus {
	if bus == nil {
		lock.Lock()
		defer lock.Unlock()
		if bus == nil {
			bus = NewEventBus()
		}
	}
	return bus
}

