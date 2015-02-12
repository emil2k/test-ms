package elevator

func ExampleControl() {
	c := NewControl()
	c.Add(1, 1)
	c.Add(2, 10)

	for c.Step() {
	}
	// Output:
	// todo
}
