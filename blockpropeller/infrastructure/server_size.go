package infrastructure

var (
	// ServerSizeTest represents a smaller server size, useful for test purposes.
	ServerSizeTest = NewServerSize("test")
	// ServerSizeProd represents a bigger server size, useful for production purposes.
	ServerSizeProd = NewServerSize("production")

	// ValidServerSizes that are recognized by BlockPropeller.
	ValidServerSizes = []ServerSize{
		ServerSizeTest,
		ServerSizeProd,
	}
)

// ServerSize defines a valid Server size.
type ServerSize string

// NewServerSize returns a new ServerSize instance.
func NewServerSize(size string) ServerSize {
	return ServerSize(size)
}

// IsValid checks whether the ServerSize is one of recognized values.
func (size ServerSize) IsValid() bool {
	for _, valid := range ValidServerSizes {
		if size == valid {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (size ServerSize) String() string {
	return string(size)
}
