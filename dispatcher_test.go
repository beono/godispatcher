package godispatcher

import (
	"errors"
	"testing"
)

func TestDispatcher_NoListeners(t *testing.T) {

	dispatcher := New()
	dispatcher.On("some_event", func(event *Event) error {
		t.Error("this listener must not be called")
		return nil
	})

	err := dispatcher.Emit("missing_event", "test")
	if err != nil {
		t.Errorf("unexpected error %+v", err)
	}
}

func TestDispatcher_OneListener(t *testing.T) {

	eventName := "some.event"
	listenerCalls := 0

	dispatcher := New()
	dispatcher.On(eventName, func(event *Event) error {
		listenerCalls++
		return nil
	})

	err := dispatcher.Emit(eventName, "test")
	if err != nil {
		t.Errorf("unexpected error %+v", err)
	}

	if listenerCalls != 1 {
		t.Errorf("excepted to listener once. called %d", listenerCalls)
	}
}

func TestDispatcher_Error(t *testing.T) {

	eventName := "some.event"

	dispatcher := New()
	dispatcher.On(eventName, func(event *Event) error {
		return errors.New("some error")
	})
	dispatcher.On(eventName, func(event *Event) error {
		t.Error("this listener must not be called")
		return nil
	})

	err := dispatcher.Emit(eventName, "test")

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err.Error() != "some error" {
		t.Errorf("unexpected error %+v", err)
	}
}

func TestDispatcher_Complex(t *testing.T) {

	eventName := "some.event"

	listenerACalls := 0
	listenerBCalls := 0

	listenerA := func(event *Event) error {
		data, ok := event.Data.(string)
		if !ok {
			t.Error("event.Data is not a string")
		}

		if data != "test" {
			t.Errorf("event.Data want: test, got: %s", data)
		}

		listenerACalls++
		return nil
	}

	listenerB := func(event *Event) error {
		data, ok := event.Data.(string)
		if !ok {
			t.Error("event.Data is not a string")
		}

		if data != "test" {
			t.Errorf("event.Data want: test, got: %s", data)
		}

		listenerBCalls++
		return nil
	}

	dispatcher := New()
	dispatcher.On(eventName, listenerA)
	dispatcher.On(eventName, listenerB)
	dispatcher.On(eventName, func(event *Event) error {
		return errors.New("some error")
	})

	err := dispatcher.Emit(eventName, "test")

	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err.Error() != "some error" {
		t.Errorf("unexpected error %+v", err)
	}

	if listenerACalls != 1 {
		t.Errorf("excepted to listener A once. called %d", listenerACalls)
	}

	if listenerBCalls != 1 {
		t.Errorf("excepted to listener B once. called %d", listenerBCalls)
	}
}
