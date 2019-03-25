package constants

import (
	"fmt"

	"github.com/VirtusLab/render/version"
)

const (
	// Name is the application name
	Name = "render"
	// Description  is s single line description of the application
	Description = "Universal file renderer"
	// Author is the application author to display
	Author = "VirtusLab"
)

func Version() string {
	ver := version.VERSION
	if len(ver) == 0 {
		ver = "unknown"
	}
	commit := version.GITCOMMIT
	if len(commit) == 0 {
		commit = "unknown"
	}
	return fmt.Sprintf("%s-%s", ver, commit)
}
