package elevator

func ExampleControl_Pickup() {
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
