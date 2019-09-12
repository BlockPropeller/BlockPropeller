package infrastructure

import "github.com/pkg/errors"

var (
	registeredDeploymentTypes map[DeploymentType]DeploymentSpec
)

func init() {
	registeredDeploymentTypes = make(map[DeploymentType]DeploymentSpec)
}

// DeploymentSpec defines common behavior that each registered Deployment
// should provide in order to be correctly managed by the system.
type DeploymentSpec interface {
	UnmarshalConfig(map[string]string) (DeploymentConfig, error)

	HealthCheck(*Server, *Deployment) (HealthCheck, error)
}

// DeploymentConfig represents custom configuration options for each deployment.
type DeploymentConfig interface {
	MarshalMap() map[string]string
}

// RegisterDeploymentType is used to register a new type of deployment.
//
// In order for the deployment to be properly managed by the system,
// this method MUST be called before working with that Deployment type.
func RegisterDeploymentType(typ DeploymentType, spec DeploymentSpec) {
	registeredDeploymentTypes[typ] = spec
}

// getDeploymentSpec returns a DeploymentSpec for the requested type, failing if none is found.
func getDeploymentSpec(typ DeploymentType) (DeploymentSpec, error) {
	spec, ok := registeredDeploymentTypes[typ]
	if !ok {
		return nil, errors.Errorf("unknown deployment type: %s", typ)
	}

	return spec, nil
}
