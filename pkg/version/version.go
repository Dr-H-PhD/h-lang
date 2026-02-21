package version

// Version information for H-lang compiler
const (
	Major = 0
	Minor = 0
	Patch = 1
)

// String returns the version string
func String() string {
	return "0.0.001"
}

// Full returns the full version with name
func Full() string {
	return "H-lang compiler (hlc) v" + String()
}
