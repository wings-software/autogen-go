// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// ConfigureSwift configures a Swift step.
func ConfigureSwift(fsys fs.FS, pipeline *spec.Pipeline) error {
	_ = pipeline.Stages[0].Spec.(*spec.StageCI)

	// TODO

	return nil
}
