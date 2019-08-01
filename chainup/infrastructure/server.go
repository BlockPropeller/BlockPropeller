package infrastructure

import (
	"sync"
	"time"

	"chainup.dev/chainup/statemachine"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

var (
	// ErrServerNotFound is returned when a ServerRepository does not find a server to return.
	ErrServerNotFound = errors.New("server not found")
	// ErrServerAlreadyExists is returned when a Server creation is attempted with an existing ServerID.
	ErrServerAlreadyExists = errors.New("server already exists")
)

// ServerID is a unique server identifier.
type ServerID string

// NewServerID returns a new unique ServerID.
func NewServerID() ServerID {
	return ServerID(uuid.NewV4().String())
}

// String satisfies the Stringer interface.
func (id ServerID) String() string {
	return string(id)
}

// ServerBuilder helps construct provisioning servers by providing
// a fluent interface for configuring server details.
type ServerBuilder struct {
	provider ProviderType
}

// NewServerBuilder starts the process of building a server.
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

// Provider configures the provider used for provisioning the server.
func (b *ServerBuilder) Provider(provider ProviderType) *ServerBuilder {
	b.provider = provider

	return b
}

// Build assembles all the server configuration into a single server object.
func (b *ServerBuilder) Build() *Server {
	return NewServer(b.provider)
}

// Server holds all the configuration values for a single provisioning server.
type Server struct {
	ID ServerID `json:"id"`

	statemachine.Resource

	Provider ProviderType `json:"provider"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// NewServer allows you to construct a provision server in a single line.
//
// If you need a fluent interface for constructing the Server, you can use the ServerBuilder.
func NewServer(provider ProviderType) *Server {
	return &Server{
		ID: NewServerID(),

		Resource: statemachine.NewResource(StateCreated),

		Provider: provider,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ServerRepository defines an interface for storing and retrieving provisioning servers.
type ServerRepository interface {
	// Find a Server given a ServerID.
	Find(id ServerID) (*Server, error)

	// Create a new Server.
	Create(server *Server) error

	// Update an existing Server.
	Update(server *Server) error
}

// InMemoryServerRepository holds the servers inside an in-memory map.
//
// Servers are not persisted on disk and won't survive program restarts.
type InMemoryServerRepository struct {
	servers sync.Map
}

// NewInMemoryServerRepository returns a new InMemoryServerRepository instance.
func NewInMemoryServerRepository() *InMemoryServerRepository {
	return &InMemoryServerRepository{}
}

// Find a Server given a ServerID.
func (repo *InMemoryServerRepository) Find(id ServerID) (*Server, error) {
	req, ok := repo.servers.Load(id)
	if !ok {
		return nil, ErrServerNotFound
	}

	return req.(*Server), nil
}

// Create a new Server.
func (repo *InMemoryServerRepository) Create(server *Server) error {
	_, loaded := repo.servers.LoadOrStore(server.ID, server)
	if loaded {
		return ErrServerAlreadyExists
	}

	return nil
}

// Update an existing Server.
func (repo *InMemoryServerRepository) Update(server *Server) error {
	repo.servers.Store(server.ID, server)

	return nil
}
