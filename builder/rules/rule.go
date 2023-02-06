package rules

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// Rule defines a pipeline build rule.
type Rule func(workspace fs.FS, pipeline *spec.Pipeline) error
