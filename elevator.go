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

// Sort sorts the queue to traverse the least floors.
func (s *State) Sort() {
	sorted, _ := path(s.Current, []Floor{}, s.Queue)
	s.Queue = sorted
}

// path finds the shortest path using Dijkstra's algorithm sorting the visited
// nodes in the necessary order.
func path(node Floor, visited, unvisited []Floor) (oVisited, oUnvisited []Floor) {
	fmt.Println("path recursion :", node, visited, unvisited)
	if len(unvisited) == 0 {
		fmt.Println("stop recursion")
		return visited, unvisited // stop recursion
	}
	init := false
	var min Floor
	var pick Floor
	pickIndex := 0
	for i, f := range unvisited {
		if d := floorDistance(node, f); !init || d < min {
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
	s.Sort()
	// TODO remove duplicates in queue
	c.fleet[e] = s
}

// TODO direction not being used here might need to remove it
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
		// Choose closest elevator based on next location.
		if d := floorDistance(s.Current, s.Next()); pick == nil || d < min {
			min = d
			pick = &e
		}
	}
	// Order the picked elevator to go to the pickup.
	fmt.Printf("elevator %d to pickup on floor %d\n", *pick, f)
	c.Update(*pick, f)
}

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

func floorDistance(a, b Floor) Floor {
	d := a - b
	if d < 0 {
		d = -d
	}
	return d
}
