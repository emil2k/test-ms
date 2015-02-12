package elevator

import (
	"fmt"
)

type Direction int

const (
	Up Direction = iota
	Down
	Stopped
)

type Elevator int

type Floor int

type State struct {
	Current Floor
	Queue   []Floor
}

func (s State) Direction() Direction {
	x := s.Next() - s.Current
	switch {
	case x > 0:
		return Up
	case x < 0:
		return Down
	default:
		return Stopped
	}
}

func (s State) Next() Floor {
	if len(s.Queue) == 0 {
		return s.Current
	}
	return s.Queue[0]
}

type Control struct {
	step  int // count of steps
	fleet map[Elevator]State
}

func NewControl() *Control {
	return &Control{step: 1, fleet: make(map[Elevator]State)}

}

func (c *Control) Add(e Elevator, f Floor) {
	c.fleet[e] = State{f, []Floor{}}
}

func (c Control) Status() map[Elevator]State {
	return c.fleet
}

func (c *Control) Update(e Elevator, f Floor) {
	s, ok := c.fleet[e]
	if !ok {
		panic("Can't update an elevator that does not exist.")
	}
	s.Queue = append(s.Queue, f)
	// TODO sort queue
	// TODO remove duplicates in queue
	c.fleet[e] = s
}

func (c *Control) Pickup(f Floor, d Direction) {
	if len(c.fleet) == 0 {
		panic("There is no elevators operating")
	}
	var min Floor // distance in floors
	var pick *Elevator = nil
	for e, s := range c.fleet {
		// Elevator on the same floor.
		if s.Current == f {
			pick = &e
			break
		}
		// Choose closest elevator.
		d := f - s.Current
		if d < 0 {
			d = -d // absolute
		}
		if pick == nil || d < min {
			min = d
			pick = &e
		}
	}
	// Order the picked elevator to go to the pickup.
	c.Update(*pick, f)
}

func (c Control) Step() bool {
	fmt.Printf("step #%d\n", c.step)
	moved := false
	for e, s := range c.fleet {
		switch s.Direction() {
		case Up:
			fmt.Printf("\televator %d goes up to floor %d\n", e, s.Next())
			moved = true
		case Down:
			fmt.Printf("\televator %d goes down to floor %d\n", e, s.Next())
			moved = true
		case Stopped:
			fmt.Printf("\televator %d is stopped on floor %d\n", e, s.Current)
		}
		// Unqueue the stop.
		if len(s.Queue) > 0 {
			s.Queue = s.Queue[1:]
		}
	}
	c.step++
	return moved
}
