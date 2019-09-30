package infrastructure

var (
	// ServerStateRequested describes the requested server specification.
	ServerStateRequested = NewServerState("requested")

	// ServerStateProvisioning indicates that the server is being provisioned by the infrastructure provider.
	ServerStateProvisioning = NewServerState("provisioning")

	// ServerStateRunning is the final success state of a Server.
	ServerStateRunning = NewServerState("running")

	// ServerStateFailed is the terminating state representing provisioning server failure.
	ServerStateFailed = NewServerState("failed")

	// ValidServerStates that are recognized by ChainUP.
	ValidServerStates = []ServerState{
		ServerStateRequested,
		ServerStateProvisioning,
		ServerStateRunning,
		ServerStateFailed,
	}
)

// ServerState defines a valid Server state.
type ServerState string

// NewServerState returns a new ServerState instance.
func NewServerState(state string) ServerState {
	return ServerState(state)
}

// IsValid checks whether the ServerState is one of recognized values.
func (state ServerState) IsValid() bool {
	for _, valid := range ValidServerStates {
		if state == valid {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (state ServerState) String() string {
	return string(state)
}
