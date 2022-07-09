//go:build buildtagdoesnotexist
// +build buildtagdoesnotexist

package build

// This file is a workaround to make sure go mod keeps around the brond and fuzz
// dependencies in the go.sum file that we only use during certain tasks (such
// as integration tests or fuzzing) or only for certain operating systems. For
// example, the specific brond import makes sure the indirect dependency
// github.com/brsuite/winsvc is kept in the go.sum file. Because of the build
// tag, this dependency never ends up in the final broln binary.
import (
	_ "github.com/brsuite/brond"
	_ "github.com/dvyukov/go-fuzz/go-fuzz-dep"
)
