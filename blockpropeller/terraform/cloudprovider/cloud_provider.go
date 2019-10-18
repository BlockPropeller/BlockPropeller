package cloudprovider

import (
	"blockpropeller.dev/blockpropeller/infrastructure"
	"blockpropeller.dev/blockpropeller/terraform"
	"github.com/pkg/errors"
)

var (
	//@TODO: This pattern diverges from the wire dependency injection. We might want to move this to wire style config.
	registeredCloudProviders map[infrastructure.ProviderType]CloudProvider
)

func init() {
	registeredCloudProviders = make(map[infrastructure.ProviderType]CloudProvider)
}

// CloudProvider is an abstraction over different cloud infrastructure providers,
// providing a common interface of provisioning infrastructure over all of them
// under a single interface.
type CloudProvider interface {
	Register(workspace *terraform.Workspace, settings *infrastructure.ProviderSettings) error
	AddServer(workspace *terraform.Workspace, srv *infrastructure.Server) error
}

// RegisterProvider is used to register a new type of cloud provider.
//
// In order for the provider to be properly managed by the system,
// this method MUST be called before working with that CloudProvider type.
func RegisterProvider(typ infrastructure.ProviderType, provider CloudProvider) {
	registeredCloudProviders[typ] = provider
}

// GetProvider returns a CloudProvider for the requested type, failing if none is found.
func GetProvider(typ infrastructure.ProviderType) (CloudProvider, error) {
	spec, ok := registeredCloudProviders[typ]
	if !ok {
		return nil, errors.Errorf("unknown provider type: %s", typ)
	}

	return spec, nil
}
