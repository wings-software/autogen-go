// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// ConfigurePython configures a Python step.
func ConfigurePython(fsys fs.FS, pipeline *spec.Pipeline) error {
	_ = pipeline.Stages[0].Spec.(*spec.StageCI)

	// TODO

	return nil
}
