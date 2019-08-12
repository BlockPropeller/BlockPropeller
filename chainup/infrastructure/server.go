package infrastructure

import (
	"net"
	"sync"
	"time"

	"chainup.dev/chainup/statemachine"
	"github.com/Pallinder/go-randomdata"
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
	name     string
	provider ProviderType
	sshKey   *SSHKey
}

// NewServerBuilder starts the process of building a server.
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

// Name configures the name used for identify the new server.
func (b *ServerBuilder) Name(name string) *ServerBuilder {
	b.name = name

	return b
}

// Provider configures the provider used for provisioning the server.
func (b *ServerBuilder) Provider(provider ProviderType) *ServerBuilder {
	b.provider = provider

	return b
}

// SSHKey configures the SSHKey used to access the server after provisioning.
func (b *ServerBuilder) SSHKey(sshKey *SSHKey) *ServerBuilder {
	b.sshKey = sshKey

	return b
}

// Build assembles all the server configuration into a single server object.
func (b *ServerBuilder) Build() (*Server, error) {
	if b.name == "" {
		b.name = randomdata.SillyName()
	}
	if b.provider == "" {
		return nil, errors.New("missing cloud provider")
	}
	if b.sshKey == nil {
		sshKey, err := GenerateNewSSHKey("ChainUP - " + b.name)
		if err != nil {
			return nil, errors.Wrap(err, "generate ssh key")
		}

		b.sshKey = sshKey
	}

	return NewServer(b.name, b.provider, b.sshKey), nil
}

// Server holds all the configuration values for a single provisioning server.
type Server struct {
	ID ServerID `json:"id"`

	statemachine.Resource

	Name string `json:"name"`

	Provider ProviderType `json:"provider"`

	SSHKey *SSHKey `json:"ssh_key"`

	IPAddress net.IP `json:"ip_address"`

	Deployments []*Deployment `json:"deployments"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// NewServer allows you to construct a provision server in a single line.
//
// If you need a fluent interface for constructing the Server, you can use the ServerBuilder.
func NewServer(name string, provider ProviderType, sshKey *SSHKey) *Server {
	return &Server{
		Resource: statemachine.NewResource(StateRequested),

		ID:       NewServerID(),
		Name:     name,
		Provider: provider,
		SSHKey:   sshKey,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AddDeployment associates a deployment with a server it is deployed on.
func (srv *Server) AddDeployment(deployment *Deployment) {
	deployment.ServerID = srv.ID
	srv.Deployments = append(srv.Deployments, deployment)
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
