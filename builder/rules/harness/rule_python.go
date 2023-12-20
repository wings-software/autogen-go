// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigurePython configures a Python step.
func ConfigurePython(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	if utils.Match(fsys, "*.py") || utils.Match(fsys, "*/*.py") {
		// add linters
		if utils.Exists(fsys, ".flake8") {
			{
				script := new(spec.StepExec)
				script.Run = "pip install flake8\nflake8 --config=.flake8 ."

				step := new(spec.Step)
				step.Name = "run linter"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
		} else if utils.Exists(fsys, ".pylintrc") {
			{
				script := new(spec.StepExec)
				script.Run = "pip install pylint\npylint --rcfile=.pylintrc ."

				step := new(spec.Step)
				step.Name = "run linter"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
		} else {
			{
				script := new(spec.StepExec)
				script.Run = "pip install flake8\nflake8 ."

				step := new(spec.Step)
				step.Name = "run linter"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
		}

		// check for the go.mod file.
		if utils.Exists(fsys, "pom.xml") {
			{
				script := new(spec.StepExec)
				script.Run = "pip install --upgrade pip setuptools"

				step := new(spec.Step)
				step.Name = "upgrade dependencies"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "mvn clean verify"
				script.Reports = []*spec.Report{
					{
						Path: []string{
							"**/*.xml",
						},
						Type: "JUnit",
					},
				}

				step := new(spec.Step)
				step.Name = "run tests and coverage"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			return nil
		} else if utils.Exists(fsys, "build.py") {
			{
				script := new(spec.StepExec)
				script.Run = "pip install --upgrade pip setuptools\npip install pybuilder"

				step := new(spec.Step)
				step.Name = "install dependencies"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "pyb"

				step := new(spec.Step)
				step.Name = "run tests and coverage"
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
			return nil
		} else if utils.Exists(fsys, "tox.ini") {
			{
				script := new(spec.StepExec)
				script.Run = "pip install tox"

				step := new(spec.Step)
				step.Name = "install dependencies"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				script := new(spec.StepExec)
				script.Run = "tox"

				step := new(spec.Step)
				step.Name = "run tests"
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

			{
				script := new(spec.StepExec)
				script.Run = "tox -e coverage"

				step := new(spec.Step)
				step.Name = "get coverage"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
			return nil
		} else if utils.Exists(fsys, "WORKSPACE") {
			{
				script := new(spec.StepExec)
				script.Run = "bazel test :*"

				step := new(spec.Step)
				step.Name = "run tests"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}
			return nil
		} else {
			{
				script := new(spec.StepExec)
				script.Run = "python3 -m venv .venv\n. .venv/bin/activate\npython3 -m pip install -r requirements.txt\npython3 -m pip install -e ."

				step := new(spec.Step)
				step.Name = "setup virtual environment"
				step.Type = "script"
				step.Spec = script

				stage.Steps = append(stage.Steps, step)
			}

			{
				test := new(spec.StepTest)
				test.With = map[string]interface{}{
					"language":           "python",
					"args":               "--junitxml=out_report.xml",
					"run_selected_tests": "true",
					"pre_command":        ". .venv/bin/activate",
				}
				test.Uses = "pytest"
				test.Reports = []*spec.Report{
					{
						Path: []string{
							"**/*.xml",
						},
					},
				}
				test.Envs = map[string]string{
					"PYTHONPATH": "/harness",
				}
				test.Splitting = &spec.Splitting{
					Enabled:  true,
					Strategy: "class_timing",
				}

				step := new(spec.Step)
				step.Name = "Run Tests"
				step.Type = "test"
				step.Spec = test
				step.Strategy = &spec.Strategy{
					Type: "for",
					Spec: &spec.For{
						Iterations: 2,
					},
				}
				stage.Steps = append(stage.Steps, step)
			}
		}
	}
	return nil
}
