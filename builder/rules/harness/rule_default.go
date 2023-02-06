// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureDefault configures a default step if the system
// is unable to automatically add any language-specific steps.
func ConfigureDefault(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// ignore if stage already contains steps
	if len(stage.Steps) > 0 {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	var image string
	if utils.IsContainerRuntime(pipeline) {
		image = "alpine"
	}

	// add dummy hello world step
	stage.Steps = append(stage.Steps, utils.CreateScriptStep(image, "echo", "echo hello world"))
	return nil
}
