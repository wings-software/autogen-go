// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureGo configures a Go step.
func ConfigureGo(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check for the go.mod file.
	if !utils.Exists(fsys, "go.mod") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	useImage := utils.IsContainerRuntime(pipeline)

	// add the go build step
	{
		script := new(spec.StepExec)
		script.Run = "go build"

		if useImage {
			script.Image = "golang"
		}

		step := new(spec.Step)
		step.Name = "go_build"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add the go test step
	{
		script := new(spec.StepExec)
		script.Run = "go test -v ./..."

		if useImage {
			script.Image = "golang"
		}

		step := new(spec.Step)
		step.Name = "go_test"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	return nil
}
