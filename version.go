package blkinfo

var (
	version  = "0.0.1"
	revision = "devel"
)

func Version() string {
	return version
}

func Revision() string {
	return revision
}
