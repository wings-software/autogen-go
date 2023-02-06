// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package drone

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/drone/spec/dist/go/convert/drone"
	"github.com/wings-software/autogen-go/utils"
)

// FromDrone configures a pipeline from a Drone yaml. If a Drone
// yaml exists in the repository we can use to more accurately
// generate the pipeline
func FromDrone(fsys fs.FS, pipeline *spec.Pipeline) error {

	// check if a Drone yaml exists
	if !utils.Exists(fsys, ".drone.yml") {
		return nil
	}

	// get the original drone yaml
	old, err := utils.Read(fsys, ".drone.yml")
	if err != nil {
		return err
	}

	// migrate to the new yaml format
	new, err := drone.FromBytes(old)
	if err != nil {
		return err
	}

	// unmarshal the new yaml
	pipeline_, err := spec.ParseBytes(new)
	if err != nil {
		return err
	}

	// copy into the pipeline stuct
	pipeline.Name = pipeline_.Name
	pipeline.Inputs = pipeline_.Inputs
	pipeline.Registry = pipeline_.Registry
	pipeline.Stages = pipeline_.Stages
	pipeline.Version = pipeline_.Version

	// override a subset of values (short term hack)
	for _, stage := range pipeline.Stages {
		stage_, ok := stage.Spec.(*spec.StageCI)
		if ok {
			// set the default platform if nil
			if stage_.Platform == nil {
				stage_.Platform = &spec.Platform{
					Os:   spec.OSLinux,
					Arch: spec.ArchAmd64,
				}
			}
			// reset the runtime (for now)
			stage_.Runtime = nil
		}
	}

	return utils.SkipAll // Skip all remaining rules
}
