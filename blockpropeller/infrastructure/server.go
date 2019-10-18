package infrastructure

import (
	"context"
	"net"
	"sync"
	"time"

	"blockpropeller.dev/blockpropeller/account"
	"blockpropeller.dev/blockpropeller/terraform"
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

var (
	// NilServerID is the value of a nil ServerID.
	NilServerID ServerID
)

// ServerID is a unique server identifier.
type ServerID string

// NewServerID returns a new unique ServerID.
func NewServerID() ServerID {
	return ServerID(uuid.NewV4().String())
}

// ServerIDFromString creates a ServerID from a string.
func ServerIDFromString(id string) ServerID {
	return ServerID(id)
}

// String satisfies the Stringer interface.
func (id ServerID) String() string {
	return string(id)
}

// ServerBuilder helps construct provisioning servers by providing
// a fluent interface for configuring server details.
type ServerBuilder struct {
	accountID account.ID
	name      string
	provider  ProviderType
	size      ServerSize
	sshKey    *SSHKey
}

// NewServerBuilder starts the process of building a server.
func NewServerBuilder(accountID account.ID) *ServerBuilder {
	return &ServerBuilder{
		accountID: accountID,
	}
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

// Size configures the size for provisioning the server.
func (b *ServerBuilder) Size(size ServerSize) *ServerBuilder {
	b.size = size

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
	if b.size == "" {
		b.size = ServerSizeTest
	}
	if b.sshKey == nil {
		sshKey, err := GenerateNewSSHKey("BlockPropeller - " + b.name)
		if err != nil {
			return nil, errors.Wrap(err, "generate ssh key")
		}

		b.sshKey = sshKey
	}

	return NewServer(b.accountID, b.name, b.provider, b.size, b.sshKey), nil
}

// Server holds all the configuration values for a single provisioning server.
type Server struct {
	ID        ServerID   `json:"id" gorm:"type:varchar(36) not null"`
	AccountID account.ID `json:"-" gorm:"type:varchar(36) not null references accounts(id)" `

	State ServerState `json:"state" gorm:"type:varchar(20) not null"`

	Name string `json:"name" gorm:"type:varchar(255) not null"`

	Provider ProviderType `json:"provider" gorm:"type:varchar(100) not null"`
	Size     ServerSize   `json:"size" gorm:"type:varchar(100) not null"`

	SSHKey *SSHKey `json:"ssh_key" gorm:"embedded;embedded_prefix:ssh_key_"`

	IPAddress net.IP `json:"ip_address,omitempty" gorm:"type:varchar(255)"`

	Deployments []*Deployment `json:"deployments,omitempty"`

	WorkspaceSnapshot *terraform.WorkspaceSnapshot `json:"-" gorm:"embedded;embedded_prefix:terraform_"`

	CreatedAt   time.Time  `json:"created_at" gorm:"type:datetime not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"type:datetime not null;default:CURRENT_TIMESTAMP"`
	CompletedAt *time.Time `json:"completed_at,omitempty" gorm:"type:datetime"`
	DeletedAt   *time.Time `json:"-" gorm:"type:datetime"`
}

// NewServer allows you to construct a provision server in a single line.
//
// If you need a fluent interface for constructing the Server, you can use the ServerBuilder.
func NewServer(accountID account.ID, name string, provider ProviderType, size ServerSize, sshKey *SSHKey) *Server {
	return &Server{
		ID:        NewServerID(),
		AccountID: accountID,

		State: ServerStateRequested,

		Name:     name,
		Provider: provider,
		Size:     size,

		SSHKey: sshKey,

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// AfterFind GORM Hook.
func (srv *Server) AfterFind() error {
	if srv.WorkspaceSnapshot != nil && srv.WorkspaceSnapshot.WorkspacePath == "" {
		srv.WorkspaceSnapshot = nil
	}

	return nil
}

// AddDeployment associates a deployment with a server it is deployed on.
func (srv *Server) AddDeployment(deployment *Deployment) {
	deployment.ServerID = srv.ID
	srv.Deployments = append(srv.Deployments, deployment)
}

// ServerRepository defines an interface for storing and retrieving provisioning servers.
type ServerRepository interface {
	// Find a Server given a ServerID.
	Find(ctx context.Context, id ServerID) (*Server, error)

	// List all servers for a particular Account.
	List(ctx context.Context, accountID account.ID) ([]*Server, error)

	// Create a new Server.
	Create(ctx context.Context, server *Server) error

	// Update an existing Server.
	Update(ctx context.Context, server *Server) error

	// Delete an existing Server
	Delete(ctx context.Context, server *Server) error
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
func (repo *InMemoryServerRepository) Find(ctx context.Context, id ServerID) (*Server, error) {
	req, ok := repo.servers.Load(id)
	if !ok {
		return nil, ErrServerNotFound
	}

	return req.(*Server), nil
}

// List all servers for a particular Account.
func (repo *InMemoryServerRepository) List(ctx context.Context, accountID account.ID) ([]*Server, error) {
	var servers []*Server

	repo.servers.Range(func(key, v interface{}) bool {
		srv := v.(*Server)
		if srv.AccountID != accountID {
			return true
		}

		servers = append(servers, srv)

		return true
	})

	return servers, nil
}

// Create a new Server.
func (repo *InMemoryServerRepository) Create(ctx context.Context, server *Server) error {
	_, loaded := repo.servers.LoadOrStore(server.ID, server)
	if loaded {
		return ErrServerAlreadyExists
	}

	return nil
}

// Update an existing Server.
func (repo *InMemoryServerRepository) Update(ctx context.Context, server *Server) error {
	repo.servers.Store(server.ID, server)

	return nil
}

// Delete an existing Server
func (repo *InMemoryServerRepository) Delete(ctx context.Context, srv *Server) error {
	repo.servers.Delete(srv.ID)

	return nil
}
