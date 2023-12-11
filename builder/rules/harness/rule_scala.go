// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"
	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureScala configures a Scala step.
func ConfigureScala(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check for the build.gradle file.
	if !utils.Exists(fsys, "build.gradle") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	useImage := utils.IsContainerRuntime(pipeline)

	// add the gradle build step
	{
		script := new(spec.StepExec)
		script.Run = "gradle build "

		if useImage {
			script.Image = "scala"
		}

		step := new(spec.Step)
		step.Name = "scala_build"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add the scala test step
	{
		script := new(spec.StepExec)
		script.Run = `./gradlew test
		xml_file=$(find /harness/build/test-results/test -type f -name "*.xml" | head -n 1)

		if [ -n "$xml_file" ]; then
			cp "$xml_file" /harness/reports.xml
			echo "XML file copied to /harness/reports.xml"
		else
			echo "No XML file found in /path"
		fi`
			
		script.Reports = append(script.Reports, &spec.Report{
            Type: "junit",
            Path: []string{"/harness/reports.xml"},
        })

		if useImage {
			script.Image = "scala"
		}

		step := new(spec.Step)
		step.Name = "scala_test"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	return nil
}
