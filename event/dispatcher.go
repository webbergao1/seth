package event

import (
	"reflect"
	"sort"
	"sync"
)

// Listener type for defining functions as listeners
type Listener struct {

	//Callable call function
	Callable func(e Event)

	//Priority priority for listener
	Priority int
}

// ListenersByPriority Listeners By Priority
type ListenersByPriority []Listener

func (l ListenersByPriority) Len() int           { return len(l) }
func (l ListenersByPriority) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l ListenersByPriority) Less(i, j int) bool { return l[i].Priority < l[j].Priority }

// Dispatcher interface defines the event dispatcher behavior
type Dispatcher interface {

	// Dispatch dispatches the event and returns it after all listeners do their jobs.
	Dispatch(e Event) Event

	// AddListener registers a listener for given event type.
	AddListener(eventType TypeEvent, listener Listener)

	// AddListenerExecOnce registers a listener to be executed only once.
	AddListenerExecOnce(eventType TypeEvent, listener Listener)

	// RemoveListener removes the registered event listener for given event type.
	RemoveListener(eventType TypeEvent, listener Listener)

	// RemoveAll removes all listeners for given type.
	RemoveAll(eventType TypeEvent)

	// HasListeners returns true if any listener for given event type
	HasListeners(eventType TypeEvent) bool
}

type listenersCollection []Listener

// eventDispatcher The EventDispatcher type is the default implementation of the DispatcherInterface
type eventDispatcher struct {
	sync.RWMutex
	listeners map[TypeEvent]listenersCollection
	sorted    map[TypeEvent]listenersCollection
}

var _instance Dispatcher

//SharedDispatcher singleton dispatcher
func SharedDispatcher() Dispatcher {
	if _instance == nil {
		_instance = NewEventDispatcher()
	}
	return _instance
}

// NewEventDispatcher creates a new instance of event dispatcher
func NewEventDispatcher() Dispatcher {
	return &eventDispatcher{
		listeners: make(map[TypeEvent]listenersCollection),
		sorted:    make(map[TypeEvent]listenersCollection),
	}
}

// AddListener registers a listener for given event type.
func (d *eventDispatcher) AddListener(eventType TypeEvent, listener Listener) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()
	d.listeners[eventType] = append(d.listeners[eventType], listener)
}

// AddListenerExecOnce registers a listener to be executed only once.
func (d *eventDispatcher) AddListenerExecOnce(eventType TypeEvent, listener Listener) {

	nl := executeRemove(d, eventType, listener) // Create a new listener that removes given listener after calling it
	d.AddListener(eventType, nl)

}

func executeRemove(d *eventDispatcher, t TypeEvent, l Listener) Listener {
	var nl Listener
	nl = Listener{
		Callable: func(e Event) {
			l.Callable(e)
			d.RWMutex.RUnlock() // The dispatcher is locked in the Dispatch method, need to unlock it
			d.RemoveListener(t, nl)
			d.RWMutex.RLock()
		},
		Priority: l.Priority,
	}

	return nl
}

// RemoveListener removes the registered event listener for given event name.
func (d *eventDispatcher) RemoveListener(eventType TypeEvent, listener Listener) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	p := reflect.ValueOf(listener.Callable).Pointer()

	listeners := d.listeners[eventType]
	for i, l := range listeners {
		lp := reflect.ValueOf(l.Callable).Pointer()
		if lp == p {
			d.listeners[eventType] = append(listeners[:i], listeners[i+1:]...)
		}
	}
}

// RemoveAll removes all listeners for given type.
func (d *eventDispatcher) RemoveAll(eventType TypeEvent) {
	d.RWMutex.Lock()
	defer d.RWMutex.Unlock()

	_, ok := d.listeners[eventType]
	if ok != false {
		delete(d.listeners, eventType)
	}
}

// HasListeners returns true if any listener for given event type
func (d *eventDispatcher) HasListeners(eventType TypeEvent) bool {
	listeners, ok := d.listeners[eventType]
	if ok == false {
		return false
	}

	return len(listeners) != 0
}

// Dispatch dispatches the event and returns it after all listeners do their jobs
func (d *eventDispatcher) Dispatch(e Event) Event {
	d.RWMutex.RLock()
	defer d.RWMutex.RUnlock()

	if !d.HasListeners(e.EventType()) {
		return e
	}

	doDispatch(d.getListeners(e.EventType()), e)

	return e
}

func doDispatch(listeners []Listener, event Event) {
	for k := range listeners {
		listeners[k].Callable(event)
		if event.IsPropagationStopped() {
			break
		}
	}

	return
}

func (d *eventDispatcher) getListeners(eventType TypeEvent) []Listener {
	if nil == d.sorted[eventType] || (len(d.sorted[eventType]) != len(d.listeners[eventType])) {
		d.sortListeners(eventType)
	}
	return d.sorted[eventType]
}

func (d *eventDispatcher) sortListeners(eventType TypeEvent) {
	sort.Sort(ListenersByPriority(d.listeners[eventType]))
	if d.sorted[eventType] == nil {
		d.sorted[eventType] = make([]Listener, 1)
	}
	d.sorted[eventType] = d.listeners[eventType]
}
