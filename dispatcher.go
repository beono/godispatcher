package godispatcher

import (
	"sync"
	"sort"
)

const DEFAULT_PRIORITY = 1000

// Event is used to transfer event-related data.
// It is passed to listeners when Emit() is called
type Event struct {
	Name string
	Data interface{}
}

// Listener is a type for listeners
// remove type Listener func(event *Event) error
type Listener struct {
	Callback func(event *Event) error
	Priority int
}

type Listeners []Listener

// Dispatcher (a.k.a event emitter, dispatcher) stores listener and notifies them when an event emitted
type Dispatcher struct {
	mutex     *sync.RWMutex
	listeners map[string]Listeners
}

func (self Listeners) Len() int {
	return len(self)
}

func (self Listeners) Less(i, j int) bool {
	return self[i].Priority < self[j].Priority
}

func (self Listeners) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

// New returns new Dispatcher
func New() Dispatcher {
	return Dispatcher{
		mutex:     &sync.RWMutex{},
		listeners: make(map[string]Listeners),
	}
}

// On adds new listener.
// listener is a callback function that will be called when event emits
// remove func (o Dispatcher) On(event string, listener Listener) {
// 	o.mutex.Lock()
// 	defer o.mutex.Unlock()
//
// 	o.listeners[event] = append(o.listeners[event], listener)
// }
func (o Dispatcher) On(event string, v interface{}) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	switch v.(type) {
	case Listener:
		if l, ok := v.(Listener); ok {
			o.listeners[event] = append(o.listeners[event], l)
		}
	case func(event *Event) error:
		if l, ok := v.(func(event *Event) error); ok {
			o.listeners[event] = append(o.listeners[event], Listener{
				Callback: l,
				Priority: DEFAULT_PRIORITY,
			})
		}
	}
}

// Emit notifies listeners about the event
func (o Dispatcher) Emit(eventName string, data interface{}) error {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	listeners, ok := o.listeners[eventName]
	if !ok {
		return nil
	}

	// add sort
	sort.Sort(listeners)

	for i := range listeners {
		event := &Event{
			Name: eventName,
			Data: data,
		}
		if err := listeners[i].Callback(event); err != nil {
			return err
		}
	}
	return nil
}
