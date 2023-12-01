// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureRust configures a Rust step.
func ConfigureRust(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// TODO

	// check for the cargo.toml file
	if !utils.Exists(fsys, "Cargo.toml") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	useImage := utils.IsContainerRuntime(pipeline)

	// add cargo build step
	{
		script := new(spec.StepExec)
		script.Run = "cargo build"

		if useImage {
			script.Image = "rust"
		}

		step := new(spec.Step)
		step.Name = "cargo_build"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add the cargo test step
	{
		script := new(spec.StepExec)
		script.Run = "cargo test "

		if useImage {
			script.Image = "rust"
		}

		step := new(spec.Step)
		step.Name = "cargo_test"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add docker build and push step
	{
		step := new(spec.Step)
		step.Name = "docker build push"
		step.Type = "BuildAndPushDockerRegistry"
		script := map[string]interface{}{
			"tags": "latest",
			"repo": "harness-test",
		}
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	return nil
}
