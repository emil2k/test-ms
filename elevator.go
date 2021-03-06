// Package elevator is a control system for a fleet of elevators.
// The goal of the system is to optimize the total path traveled by the
// elevator fleet. It uses an implementation of the Dijkstra's algorithm to
// determine the shortest path and sort the queue.
// It can queue multiple floor destinations per elevator, and makes no
// assumption based on direction during pickup - so only one button.
package elevator

import (
	"fmt"
)

// Elevator represents an elevator.
type Elevator int

// Floor represents a floor or a difference between floors.
type Floor int

// Direction is an enum value representing the direction of the elevator at any
// given time.
type Direction int

const (
	Up Direction = iota
	Down
	Stopped
)

// State represents the current state of the elevator, including the current
// floor and the sorted queue of floors it needs to visit next.
type State struct {
	Current Floor
	Queue   []Floor
}

// Direction returns the direction the elevator is moving next.
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

// Next returns the next floor the elevator will visit.
func (s State) Next() Floor {
	if len(s.Queue) == 0 {
		return s.Current
	}
	return s.Queue[0]
}

// Enqueue appends a floor to the end of the queue if it is not already present.
func (s *State) Enqueue(f Floor) {
	for _, v := range s.Queue {
		if v == f {
			return
		}
	}
	s.Queue = append(s.Queue, f)
}

// Sort sorts the queue to traverse the least floors.
func (s *State) Sort() {
	sorted, _ := path(s.Current, []Floor{}, s.Queue)
	s.Queue = sorted
}

// Total is the total distance in floors that needs to be traveled to satisfy
// current queue.
func (s *State) Total() Floor {
	var total Floor = 0
	last := s.Current
	for _, v := range s.Queue {
		total += distance(last, v)
		last = v
	}
	return total
}

// path finds the shortest path using Dijkstra's algorithm sorting the visited
// nodes in the necessary order.
func path(node Floor, visited, unvisited []Floor) (oVisited, oUnvisited []Floor) {
	if len(unvisited) == 0 {
		return visited, unvisited // stop recursion
	}
	init := false
	var min Floor
	var pick Floor
	pickIndex := 0
	for i, f := range unvisited {
		if d := distance(node, f); !init || d < min {
			init = true
			min = d
			pick = f
			pickIndex = i
		}
	}
	oVisited = append(visited, pick)
	oUnvisited = append(unvisited[:pickIndex], unvisited[pickIndex+1:]...)
	return path(pick, oVisited, oUnvisited)
}

// Control controls a fleet of elevators and determines their movements to
// minimize their total traversed path.
type Control struct {
	step  int // count of steps
	fleet map[Elevator]State
}

func NewControl() *Control {
	return &Control{step: 1, fleet: make(map[Elevator]State)}

}

// Add adds an elevator at the given initial floor.
func (c *Control) Add(e Elevator, f Floor) {
	c.fleet[e] = State{f, []Floor{}}
}

func (c Control) Status() map[Elevator]State {
	return c.fleet
}

// Update directs and elevator to visit a certain floor, by appending it to its
// queue and sorting it.
func (c *Control) Update(e Elevator, f Floor) {
	s, ok := c.fleet[e]
	if !ok {
		panic("Can't update an elevator that does not exist.")
	}
	s.Enqueue(f)
	s.Sort()
	c.fleet[e] = s
}

// Pickup tells the control system that someone needs to be picked up at the
// specified floor. The control system determines the elevator to execute the
// pickup.
func (c *Control) Pickup(f Floor) {
	if len(c.fleet) == 0 {
		panic("There is no elevators operating")
	}
	init := false
	var min Floor // distance in floors
	var pick Elevator
	for e, s := range c.fleet {
		// Determine which elevator would have the shortest path with
		// this pickup in its queue.
		s.Enqueue(f)
		s.Sort()
		x := s.Total()
		if !init || x < min {
			init = true
			min = x
			pick = e
		}
	}
	// Order the picked elevator to go to the pickup.
	fmt.Printf("elevator %d to pickup on floor %d\n", pick, f)
	c.Update(pick, f)
}

// Step ouputs each elevators action for the current step and preps queues for
// the next step.
func (c *Control) Step() bool {
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
		// Move that elevator.
		if len(s.Queue) > 0 {
			s.Current = s.Queue[0]
			s.Queue = s.Queue[1:]
			c.fleet[e] = s
		}
	}
	c.step++
	return moved
}

// distance provides the distance in floors between to floors.
func distance(a, b Floor) Floor {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d
}
