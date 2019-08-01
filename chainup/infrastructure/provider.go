package infrastructure

var (
	// ProviderDigitalOcean is the ProviderType for DigitalOcean cloud provider.
	ProviderDigitalOcean ProviderType = "digitalocean"

	// ValidProviders that are recognized by ChainUP.
	ValidProviders = []ProviderType{ProviderDigitalOcean}
)

// ProviderType that is able to provision new infrastructure.
type ProviderType string

// NewProviderType returns a new ProviderType instance.
func NewProviderType(provider string) ProviderType {
	return ProviderType(provider)
}

// IsValid checks if the provider is one of the ChainUP recognized providers.
func (pt ProviderType) IsValid() bool {
	for _, valid := range ValidProviders {
		if pt == valid {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (pt ProviderType) String() string {
	return string(pt)
}

//@TODO: Add a method for registering providers similarly to how sql package registers its drivers.
