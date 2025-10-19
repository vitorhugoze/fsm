package fsm

import (
	"math"
)

type Fsm[T ~int] struct {
	currentState T
	states       []State[T]
}

func NewFsm[T ~int](initialState T) Fsm[T] {
	return Fsm[T]{
		currentState: initialState,
	}
}

func (fsm Fsm[T]) AddState(state State[T]) Fsm[T] {
	fsm.states = append(fsm.states, state)
	return fsm
}

func (fsm *Fsm[T]) CurrentState() T {
	return fsm.currentState
}

func (fsm *Fsm[T]) Event(event string) {
	for i, s := range fsm.states {
		if s.state != fsm.currentState {
			continue
		}

		state := fsm.states[i]
		if state.events[event] == nil {
			return
		}

		next, err := state.events[event]()
		if err != nil {
			fsm.fallback()
			return
		}

		fsm.changeState(next)
		return
	}
}

func (fsm *Fsm[T]) changeState(toState T) {
	err := fsm.triggerCallback(EXIT, fsm.currentState)
	if err != nil {
		fsm.fallback()
		return
	}

	fsm.currentState = toState
	err = fsm.triggerCallback(ENTER, toState)
	if err != nil {
		fsm.fallback()
		return
	}
}

func (fsm *Fsm[T]) triggerCallback(phase EnCallbackPhase, state T) error {
	for i, s := range fsm.states {
		if s.state != fsm.currentState {
			continue
		}

		state := fsm.states[i]
		if state.callbacks[phase] == nil {
			return nil
		}

		err := state.callbacks[phase]()
		return err
	}

	return nil
}

func (fsm *Fsm[T]) fallback() {
	for i, s := range fsm.states {
		if s.state != fsm.currentState {
			continue
		}

		state := fsm.states[i]
		if state.fallback != math.MinInt {
			fsm.currentState = state.fallback
		}

		return
	}
}

type EnCallbackPhase int

const (
	ENTER EnCallbackPhase = iota + 1
	EXIT
)

type State[T ~int] struct {
	state     T
	fallback  T
	events    map[string]func() (T, error)
	callbacks map[EnCallbackPhase]func() error
}

// Initial fallback is set to math.MinInt, so 0 can be a valid fallback state
func NewState[T ~int](state T) State[T] {
	return State[T]{
		state:     state,
		fallback:  math.MinInt,
		events:    make(map[string]func() (T, error)),
		callbacks: make(map[EnCallbackPhase]func() error),
	}
}

/*
Event is the main entity of a State.
When a event is executed the next step is returned.
If an error occurs, the next step will be the fallback if it is set.
In case of an error ocurring without a fallback being set, the Fsm will continue at the current State.
*/
func (state State[T]) AddEvent(eventName string, event func() (T, error)) State[T] {
	state.events[eventName] = event
	return state
}

/*
Callback can be set to execute functions when entering and exiting the state.
Callbacks just change the Fsm state if an error occurs and a valid fallback is set.
*/
func (state State[T]) AddCallback(phase EnCallbackPhase, callback func() error) State[T] {
	state.callbacks[phase] = callback
	return state
}

func (state State[T]) WithFallback(fallback T) State[T] {
	state.fallback = fallback
	return state
}

type MyStates int

const (
	Start MyStates = iota
	Plan1
	Plan2
	Plan3
	Confirm
)
