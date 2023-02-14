// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package harness

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/wings-software/autogen-go/utils"
)

// ConfigurePlatform configures the platform
func ConfigurePlatform(fsys fs.FS, pipeline *spec.Pipeline) error {
	stage := pipeline.Stages[0].Spec.(*spec.StageCI)

	// assume the os and architecture are linux / amd64
	// stage.Platform = &spec.Platform{
	// 	Os:   spec.OSLinux,
	// 	Arch: spec.ArchAmd64,
	// }

	// unless this is an xcode project
	if isXcode(fsys) {
		stage.Platform.Os = spec.OSMacos
		stage.Platform.Arch = spec.ArchArm64
	}

	return nil
}

// helper function returns true if the project has an xcode directory
// in the root or a subdirectory of the repository.
func isXcode(workspace fs.FS) bool {
	return utils.Match(workspace, "*.xcodeproj") || utils.Match(workspace, "*/*.xcodeproj")
}
