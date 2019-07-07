package blkinfo

var (
	version  = "0.0.1"
	revision = "devel"
)

// Version returns the version.
func Version() string {
	return version
}

// Revision returns the revision.
func Revision() string {
	return revision
}
