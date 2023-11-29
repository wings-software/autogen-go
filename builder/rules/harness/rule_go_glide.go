// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"fmt"
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureGo configures a Go step.
func ConfigureGoGlide(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check for the go.mod file.
	if !utils.Exists(fsys, "vendor") && !utils.Exists(fsys, "glide.yaml") && !utils.Exists(fsys, "glide.lock") {
		return nil
	}

	fmt.Println("Found vendor folder")
	// check if we should use a container-based
	// execution environment.
	useImage := utils.IsContainerRuntime(pipeline)

	// add the glide prerequisite step
	{
		script := new(spec.StepExec)
		script.Run = `export GOBIN=/home/harness/go/bin
		go install github.com/jstemmer/go-junit-report/v2@latest
		mkdir /home/harness/go/src
		
		echo cd to gopath
		cd /home/harness/go
		ls
		
		export PATH=/home/harness/go/bin:$PATH
		echo $PATH
		
		
		# GLIDE_VERSION="v0.13.3"
		GLIDE_URL="https://github.com/Masterminds/glide/releases/download/v0.13.3/glide-v0.13.3-linux-386.zip"
		TMP_DIR="/tmp/glide_tmp"
		
		# Create a temporary directory
		mkdir -p "$TMP_DIR"
		
		# Download Glide binary
		wget "$GLIDE_URL" -O "$TMP_DIR/glide.zip"
		
		# Extract the downloaded ZIP file
		unzip "$TMP_DIR/glide.zip" -d "$TMP_DIR"
		
		# Find Glide binary in extracted folder
		# GLIDE_BINARY=$(find "$TMP_DIR" -name "glide*" -type f -executable)
		# echo $GLIDE_BINARY
		GLIDE_BINARY=$TMP_DIR/linux-386/glide
		# Move Glide binary to go/bin directory
		if [ -n "$GLIDE_BINARY" ]; then
			mv "$GLIDE_BINARY" "/home/harness/go/bin/"
			echo "Glide binary moved to /home/harness/go/bin/"
		else
			echo "Glide binary not found."
		fi
		
		# Clean up temporary directory
		rm -rf "$TMP_DIR"
		
		
		echo installed
		
		
		cd /harness
		
		export GO111MODULE=off
		
		echo glide install
		glide install
		
		export GO111MODULE=off
		go get
		echo done go get
		go build`

		if useImage {
			script.Image = "golang"
		}

		step := new(spec.Step)
		step.Name = "prerequisite"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add the go build step
	{
		script := new(spec.StepExec)
		script.Run = `export GO111MODULE=off		
		go build
		`

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
		script.Run = `export GO111MODULE=off
		go test -coverprofile=coverage.out ./...`

		if useImage {
			script.Image = "golang"
		}

		step := new(spec.Step)
		step.Name = "go_test_coverage"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	// add the go test with report step
	{
		script := new(spec.StepExec)
		script.Run = `export GOBIN=/home/harness/go/bin
		export PATH=/home/harness/go/bin:$PATH
		go install github.com/jstemmer/go-junit-report/v2@latest
		
		export GO111MODULE=off
		go test -v 2>&1 ./... | go-junit-report -set-exit-code > report.xml`

		if useImage {
			script.Image = "golang"
		}

		script.Reports = append(script.Reports, &spec.Report{
			Type: "junit",
			Path: []string{"/harness/report.xml"},
		})
		step := new(spec.Step)
		step.Name = "go_test_report"
		step.Type = "script"
		step.Spec = script

		stage.Steps = append(stage.Steps, step)
	}

	return nil
}
