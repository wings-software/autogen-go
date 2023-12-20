// Copyright 2023 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureCPP configures a CPP step.
func ConfigureCPP(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	if utils.Match(fsys, "*.cpp") || utils.Match(fsys, "*/*.cpp") || utils.Match(fsys, "*/*.cc") || utils.Match(fsys, "*.cc") {
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
				script.Run = "TEST_NODE_INDEX=<+strategy.iteration>\nTEST_NODE_TOTAL=<+strategy.iterations>\n" +
					"export TEST_FILES=`split_tests  --glob 'test/**/*.cpp' --split-by file_timing --verbose " +
					"--split-index '${TEST_NODE_INDEX}' --split-total '${TEST_NODE_TOTAL}'`"
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
		} else if utils.Exists(fsys, "CMakeLists.txt") {
			{
				script := new(spec.StepExec)
				script.Run = "sudo apt-get update\nsudo apt-get install libgtest-dev"
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
				script.Run = "cd build && make check"
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
		}
	}

	return nil
}
