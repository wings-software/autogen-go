// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"bytes"
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// ConfigureRails configures a Ruby on Rails step.
func ConfigureRails(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// check if we should use a container-based
	// execution environment.
	var image string
	if isContainerRuntime(pipeline) {
		image = "ruby"
	}

	// check for a ruby gemfile
	if exists(fsys, "Gemfile") {

		// ignore gemfiles that do not contain the
		// rails dependency
		gemfile, _ := read(fsys, "Gemfile")
		if !bytes.Contains(gemfile, []byte("'rails'")) {
			return nil
		}

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_install",
			"bundle install — jobs=3 — retry=3",
		))

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_db_create",
			"bundle exec rake db:create",
		))

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_db_migrate",
			"bundle exec rake db:migrate",
		))

		stage.Steps = append(stage.Steps, createScriptStep(image,
			"bundle_rspec",
			"bundle exec rspec",
		))
	}

	return nil
}
