package statemachine

// Middleware enables the StateMachine to be extended with additional functionality.
//@TODO: Find less verbose middleware pattern.
type Middleware interface {
	Wrap(step Step) Step
}

// MiddlewareFn is a helper function for defining middleware that does not require
// any dependencies.
type MiddlewareFn func(step Step) Step

// Wrap implements the Middleware interface.
func (fn MiddlewareFn) Wrap(step Step) Step {
	return fn(step)
}

// MiddlewareStack represents a middleware stack configured for a particular StateMachine.
type MiddlewareStack []Middleware

// Extend extends the Middleware stack with additional middleware.
func (s MiddlewareStack) Extend(middleware ...Middleware) MiddlewareStack {
	newStack := append([]Middleware{}, s...)

	return append(newStack, middleware...)
}

// Do finalizes the middleware stack and prepares a final step to be executed.
func (s MiddlewareStack) Do(step Step) Step {
	if len(s) == 0 {
		return step
	}

	for i := len(s) - 1; i >= 0; i-- {
		step = s[i].Wrap(step)
	}

	return step
}
