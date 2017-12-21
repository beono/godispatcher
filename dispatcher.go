package godispatcher

import "sync"

// Event is used to transfer event-related data.
// It is passed to listeners when Emit() is called
type Event struct {
	Name string
	Data interface{}
}

// Listener is a type for listeners
type Listener func(event *Event) error

// Dispatcher (a.k.a event emitter, dispatcher) stores listener and notifies them when an event emitted
type Dispatcher struct {
	mutex     *sync.RWMutex
	listeners map[string][]Listener
}

// New returns new Dispatcher
func New() Dispatcher {
	return Dispatcher{
		mutex:     &sync.RWMutex{},
		listeners: make(map[string][]Listener),
	}
}

// On adds new listener.
// listener is a callback function that will be called when event emits
func (o Dispatcher) On(event string, listener Listener) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.listeners[event] = append(o.listeners[event], listener)
}

// Emit notifies listeners about the event
func (o Dispatcher) Emit(eventName string, data interface{}) error {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	listeners, ok := o.listeners[eventName]
	if !ok {
		return nil
	}

	for i := range listeners {
		event := &Event{
			Name: eventName,
			Data: data,
		}
		if err := listeners[i](event); err != nil {
			return err
		}
	}
	return nil
}
