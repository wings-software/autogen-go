// Copyright 2023 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureC configures a C step.
func ConfigureC(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	if utils.Match(fsys, "*/*.c") || utils.Match(fsys, "*.c") {
		if utils.Exists(fsys, "Makefile") {
			{
				script := new(spec.StepExec)
				script.Run = "make"

				step := new(spec.Step)
				step.Name = "build make"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "make lint"

				step := new(spec.Step)
				step.Name = "make lint"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "make test"
				script.Reports = []*spec.Report{
					{
						Path: []string{
							"**/*.xml",
						},
						Type: "JUnit",
					},
				}

				step := new(spec.Step)
				step.Name = "run tests"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
		} else if utils.Match(fsys, "CMakeLists.txt") {
			{
				script := new(spec.StepExec)
				script.Run = "sudo apt-get update\nsudo apt-get install check && sudo apt-get install lcov"
				step := new(spec.Step)
				step.Name = "install and update deps"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "mkdir build && cd build && cmake .. && cmake --build ."

				step := new(spec.Step)
				step.Name = "cmake build && make"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "cd build && make test"
				script.Reports = []*spec.Report{
					{
						Path: []string{
							"**/*.xml",
						},
						Type: "JUnit",
					},
				}

				step := new(spec.Step)
				step.Name = "run tests"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
		} else if utils.Match(fsys, "Makefile.am") {
			{
				script := new(spec.StepExec)
				script.Run = "./bootstrap && ./configure"
				step := new(spec.Step)
				step.Name = "configure"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "make clean && make check && make cov"

				step := new(spec.Step)
				step.Name = "run tests and get coverage"
				step.Type = "script"
				step.Spec = script
				script.Reports = []*spec.Report{
					{
						Path: []string{
							"**/*.xml",
						},
						Type: "JUnit",
					},
				}
				stage.Steps = append(stage.Steps, step)
			}
		}
	}

	return nil
}
