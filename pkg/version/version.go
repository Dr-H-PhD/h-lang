package version

import "fmt"

// Version information for H-lang compiler
const (
	Major = 0
	Minor = 0
	Patch = 0
	Build = 334
)

// String returns the version string
func String() string {
	return fmt.Sprintf("%d.%d.%d.%03d", Major, Minor, Patch, Build)
}

// Full returns the full version with name
func Full() string {
	return "H-lang compiler (hlc) v" + String()
}
