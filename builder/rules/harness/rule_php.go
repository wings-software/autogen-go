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
func ConfigurePHP(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	if !utils.Exists(fsys, "composer.json") {
		return nil
	}

	var image string
	if utils.IsContainerRuntime(pipeline) {
		image = "php"
	}

	stage.Steps = append(stage.Steps, utils.CreateScriptWithUserDefinition(image, "composer_install", "composer install"))

	json := new(Composer)
	if err := utils.Unmarshal(fsys, "composer.json", &json); err != nil {
		return nil
	}

	if _, ok := json.Scripts["test"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptWithUserDefinitionAndReportPath(image,
			"composer_test_report",
			"composer test",
		))
	}

	if _, ok := json.Scripts["coverage"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptWithUserDefinition(image,
			"composer_test_coverage",
			"composer coverage",
		))
	}

	if _, ok := json.Scripts["lint"]; ok {
		stage.Steps = append(stage.Steps, utils.CreateScriptWithUserDefinition(image,
			"composer_lint",
			"composer lint",
		))
	}

	return nil
}

type Composer struct {
	Name    string                 `json:"name"`
	Version string                 `json:"version"`
	Scripts map[string]interface{} `json:"scripts"`
}
