// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// ConfigureDefault configures a default step if the system
// is unable to automatically add any language-specific steps.
func ConfigurePipelineVersion(fsys fs.FS, pipeline *spec.Pipeline) error {
	pipeline.Version = spec.StringorInt(1)
	return nil
}
