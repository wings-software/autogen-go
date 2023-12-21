// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureNode configures a Node step.
func ConfigureNode(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check for the package.json file.
	if !utils.Exists(fsys, "package.json") {
		return nil
	}

	// check if we should use a container-based
	// execution environment.
	var image string
	if utils.IsContainerRuntime(pipeline) {
		image = "node"
	}

	// add the npm install step
	stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
		"install dependencies",
		"npm install",
	))

	// parse the package.json file and unmarshal
	json := new(packageJson)
	if err := utils.Unmarshal(fsys, "package.json", &json); err != nil {
		return nil
	}

	// add test with junit for xml reports otherwise add well known test step
	if _, ok := json.Scripts["test:junit"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStepWithReports(image,
			"test_reports",
			"npm run test:junit",
		))
	} else if _, ok := json.Scripts["test"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"test",
			"npm run test",
		))
	}

	// add well-known jest coverage command
	if _, ok := json.Scripts["coverage"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"coverage",
			"npm run coverage",
		))
	}

	// add well-known lint command
	if _, ok := json.Scripts["lint"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"lint",
			"npm run lint",
		))
	}

	// add well-known e2e command
	if _, ok := json.Scripts["e2e"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"e2e",
			"npm run e2e",
		))
	}

	// add well-known dist command
	if _, ok := json.Scripts["dist"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"dist",
			"npm run dist",
		))
	}

	// add well-known build command
	if _, ok := json.Scripts["build"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"build",
			"npm run build",
		))
	}

	// add well-known e2e docker if infra is cloud
	if _, ok := json.Scripts["e2e:docker"]; ok && image == "" {
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"e2e_docker",
			"npm run e2e docker",
		))
	}

	return nil
}

// represents the package.json file format.
type packageJson struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Scripts map[string]interface{} `json:"scripts"`
}
