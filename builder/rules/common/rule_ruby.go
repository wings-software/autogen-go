// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package common

import (
	"bytes"
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigureRuby configures a Ruby on Rails step.
func ConfigureRuby(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check if we should use a container-based
	// execution environment.
	var image string
	if utils.IsContainerRuntime(pipeline) {
		image = "ruby"
	}

	// generate pipeline steps for rakefiles
	if utils.Exists(fsys, "Rakefile") {
		rakefile, _ := utils.Read(fsys, "Rakefile")

		// ignore ruby on rails.  we will handle rails
		// in a separate rule.
		gemfile, _ := utils.Read(fsys, "Gemfile")
		if bytes.Contains(gemfile, []byte("'rails'")) {
			return nil
		}

		// bundle install
		stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
			"bundle_install",
			"bundle install --local || bundle install",
		))

		// count of raketasks added
		var raketasks int

		// look for well known :compile command
		if bytes.Contains(rakefile, []byte(":compile")) {
			raketasks++
			stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
				"rake_compile",
				"bundle exec rake compile",
			))
		}

		// look for well known :test command
		if bytes.Contains(rakefile, []byte(":test")) {
			raketasks++
			stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
				"rake_test",
				"bundle exec rake test",
			))
		}

		// if no raketasks added run the default
		if raketasks == 0 {
			raketasks++
			stage.Steps = append(stage.Steps, utils.CreateScriptStep(image,
				"rake",
				"bundle exec rake",
			))
		}

		return nil
	}

	//
	// generate pipeline steps for ruby projects that
	// do not use rakefiles.
	//

	return nil
}
