package blkinfo

var (
	version  = "devel"
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
