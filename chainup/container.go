package chainup

// Container is an abstraction of a runnable
// that needs to be executed on a target host machine.
type Container struct {
	Name string
}

// NewContainer creates a new instance of a Container struct.
func NewContainer(name string) *Container {
	return &Container{Name: name}
}

// String satisfies the Stringer interface.
func (c Container) String() string {
	return c.Name
}
