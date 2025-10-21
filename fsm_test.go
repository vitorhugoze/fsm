package fsm

import (
	"testing"
)

type EnPhases int

const (
	Start EnPhases = iota
	Plan1
	Plan2
	Plan3
	End
)

func TestPlanSelection(t *testing.T) {
	start := NewState(Start).AddCallback(ENTER, func() error {
		println("Hi, what plan woud you like to choose ?\n 1-Standard plan\n 2-Pro plan\n 3-Enterprise plan")
		return nil
	}).AddEvent("plan1", func() (EnPhases, error) {
		return Plan1, nil
	}).AddEvent("plan2", func() (EnPhases, error) {
		return Plan2, nil
	}).AddEvent("plan3", func() (EnPhases, error) {
		return Plan3, nil
	})

	plan1 := NewState(Plan1).AddCallback(ENTER, func() error {
		println("Hi, you choose the Standard plan 1- Make payment 2- Select another plan")
		return nil
	}).AddEvent("payment", func() (EnPhases, error) {
		return End, nil
	}).AddEvent("plans", func() (EnPhases, error) {
		return Start, nil
	})

	plan2 := NewState(Plan2).AddCallback(ENTER, func() error {
		println("Hi, you choose the Pro plan 1- Make payment 2- Select another plan")
		return nil
	}).AddEvent("payment", func() (EnPhases, error) {
		return End, nil
	}).AddEvent("plans", func() (EnPhases, error) {
		return Start, nil
	})

	plan3 := NewState(Plan3).AddCallback(ENTER, func() error {
		println("Hi, you choose the Enterprise plan 1- Make payment 2- Select another plan")
		return nil
	}).AddEvent("payment", func() (EnPhases, error) {
		return End, nil
	}).AddEvent("plans", func() (EnPhases, error) {
		return Start, nil
	})

	end := NewState(End).AddCallback(ENTER, func() error {
		println("Hi, the link for your payment is http://")
		return nil
	})

	fsm := NewFsm[EnPhases]().AddState(start).AddState(plan1).AddState(plan2).AddState(plan3).AddState(end)

	fsm.SetState(Start)
	if fsm.CurrentState() != Start {
		t.Fatal("wrong state")
	}

	fsm.Event("plan1")
	if fsm.CurrentState() != Plan1 {
		t.Fatal("wrong state")
	}

	fsm.Event("plans")
	if fsm.CurrentState() != Start {
		t.Fatal("wrong state")
	}

	fsm.Event("plan2")
	if fsm.CurrentState() != Plan2 {
		t.Fatal("wrong state")
	}

	fsm.Event("plans")
	if fsm.CurrentState() != Start {
		t.Fatal("wrong state")
	}

	fsm.Event("plan3")
	if fsm.CurrentState() != Plan3 {
		t.Fatal("wrong state")
	}

	fsm.Event("payment")
	if fsm.CurrentState() != End {
		t.Fatal("wrong state")
	}
}
