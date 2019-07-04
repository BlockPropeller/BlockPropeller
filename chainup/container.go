package chainup

type Container struct {
	Name string
}

func NewContainer(name string) *Container {
	return &Container{Name: name}
}

func (c Container) String() string {
	return c.Name
}
