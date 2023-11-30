// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"

	spec "github.com/drone/spec/dist/go"
)

// SkipAll is used as a return value from Rule to indicate
// that all remaining rules are to be skipped. It is never
// returned as an error by the Builder.
var SkipAll = errors.New("skip everything and stop the pipeline generation")

// helper function returns true if the files or folders
// matching the specified pattern exist at the base path.
func Match(fsys fs.FS, pattern string) bool {
	matches, _ := fsys.(fs.GlobFS).Glob(pattern)
	return len(matches) > 0
}

// helper function returns true if the named file exists
// at the base path.
func Exists(fsys fs.FS, name string) bool {
	_, err := fsys.(fs.StatFS).Stat(name)
	return err == nil
}

// // helper function returns true if the dir exist in base path
// at the base path.
func IsDirectoryPresent(fsys fs.FS, directoryPath string) bool {
	// Use os.Stat to check if the directory exists
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		//return
	}

	fmt.Println("Current Working Directory:", currentDir)
	fmt.Println("kotlin Directory:", currentDir+directoryPath)

	_, err = os.Stat(directoryPath + directoryPath)

	// Check if there is no error and the path exists and is a directory
	if err == nil {
		fmt.Println("kotlin Directory is present:", currentDir+directoryPath)
		return true
	}

	// Check if the error is due to the directory not existing
	if os.IsNotExist(err) {
		fmt.Println("kotlin Directory is not present:", currentDir+directoryPath)
		return false
	}

	// Handle other errors (e.g., permission issues)
	fmt.Println("Error checking directory:", err)
	return false
}

// helper function reads the named file at the base path.
func Read(fsys fs.FS, name string) ([]byte, error) {
	return fsys.(fs.ReadFileFS).ReadFile(name)
}

// helper function unmarshals the named file at the base path
// into the go structure.
func Unmarshal(fsys fs.FS, name string, v interface{}) error {
	data, err := Read(fsys, name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// helper function returns true if the runtime engine is
// kubernetes or is container-based.
func IsContainerRuntime(pipeline *spec.Pipeline) bool {
	// ensure default stages already added
	if len(pipeline.Stages) == 0 {
		return false
	}
	// ensure default stage is continuous integration
	stage, ok := pipeline.Stages[0].Spec.(*spec.StageCI)
	if !ok {
		return false
	}
	// ensure runtime is not null
	if stage.Runtime == nil {
		return false
	}
	switch stage.Runtime.Type {
	case "kubernetes":
		return true
	default:
		return false
	}
}

// helper function to create a script step.
func CreateScriptStep(image, name, command string) *spec.Step {
	script := new(spec.StepExec)
	script.Run = command

	if image != "" {
		script.Image = image
	}

	step := new(spec.Step)
	step.Name = name
	step.Type = "script"
	step.Spec = script

	return step
}
