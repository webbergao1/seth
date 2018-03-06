package event

// Event interface
type Event interface {
	EventType() TypeEvent
	// IsPropagationStopped informs weather the event should
	// be further propagated or not
	IsPropagationStopped() bool

	// StopPropagation makes the event no longer propagate.
	StopPropagation()
}

// ParamsEvent is the default implementation of Event interface. Contains additional
// string parameters
type ParamsEvent struct {
	eventtype            TypeEvent
	isPropagationStopped bool
	params               map[string]interface{}
}

// NewParamsEvent is a factory for creating a basic event
func NewParamsEvent(eventtype TypeEvent) *ParamsEvent {
	p := make(map[string]interface{})
	e := ParamsEvent{eventtype, false, p} // Propagation never stopped by default
	return &e
}

// EventType returns the type of the event
func (event ParamsEvent) EventType() TypeEvent {
	return event.eventtype
}

// IsPropagationStopped informs weather the event should
// be further propagated or not
func (event ParamsEvent) IsPropagationStopped() bool {
	return event.isPropagationStopped
}

// StopPropagation sets a flag that make the event no longer propagate.
func (event *ParamsEvent) StopPropagation() {
	event.isPropagationStopped = true
}

// SetParam set a parameter for the event.
func (event *ParamsEvent) SetParam(k string, v interface{}) *ParamsEvent {
	event.params[k] = v
	return event
}

// RemoveParam deletes a param with given key.
func (event *ParamsEvent) RemoveParam(k string) *ParamsEvent {
	if event.HasParam(k) {
		delete(event.params, k)
	}
	return event
}

// HasParam defines if a param with given key exists. Returns a boolean value
func (event *ParamsEvent) HasParam(k string) bool {
	_, ok := event.params[k]
	return ok
}

// GetParam returns a parameter value for given key.
func (event *ParamsEvent) GetParam(k string) (value interface{}, ok bool) {
	v, ok := event.params[k]
	if ok == false {
		return "", false
	}
	return v, ok
}
