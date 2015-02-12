package elevator

func ExampleControl_1() {
	c := NewControl()
	c.Add(1, 5)
	c.Add(2, 1)

	c.Update(1, 10)
	c.Update(1, 1)
	c.Update(2, 4)

	c.Pickup(4, Down)
	c.Pickup(9, Up)

	for c.Step() {
	}
	// Output:
	// todo
}

func ExampleControl_2() {
	c := NewControl()
	c.Add(1, 5)
	c.Update(1, 10)
	c.Update(1, 8)
	c.Update(1, 100)
	c.Update(1, 100) // test duplicates
	c.Update(1, 5)
	for c.Step() {
	}
	// Output:
	// todo
}
