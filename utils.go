package dfquery

import (
	"runtime/debug"
	_ "unsafe"
)

func getDfVersion() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return "Unknown"
	}
	for _, dep := range bi.Deps {
		if dep.Path == "github.com/df-mc/dragonfly" {
			return dep.Version
		}
	}
	return "Unknown"
}
