package blkinfo

var (
	version  = "dev"
	revision = "none"
	builtAt  = "unknown"
)

// Version returns version.
func Version() string {
	return version
}

// Revision returns revision.
func Revision() string {
	return revision
}

// BuiltAt returns builtAt.
func BuiltAt() string {
	return builtAt
}
