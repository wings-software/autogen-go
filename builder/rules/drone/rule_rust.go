// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drone

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// ConfigureRust configures a Rust step.
func ConfigureRust(fsys fs.FS, pipeline *spec.Pipeline) error {
	_ = pipeline.Stages[0].Spec.(*spec.StageCI)

	// TODO

	return nil
}
