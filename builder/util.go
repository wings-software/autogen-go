// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package builder

import (
	"encoding/json"
	"io/fs"

	spec "github.com/drone/spec/dist/go"
)

// helper function returns true if the files or folders
// matching the specified pattern exist at the base path.
func match(fsys fs.FS, pattern string) bool {
	matches, _ := fsys.(fs.GlobFS).Glob(pattern)
	return len(matches) > 0
}

// helper function returns true if the named file exists
// at the base path.
func exists(fsys fs.FS, name string) bool {
	_, err := fsys.(fs.StatFS).Stat(name)
	return err == nil
}

// helper function reads the named file at the base path.
func read(fsys fs.FS, name string) ([]byte, error) {
	return fsys.(fs.ReadFileFS).ReadFile(name)
}

// helper function unmarshals the named file at the base path
// into the go structure.
func unmarshal(fsys fs.FS, name string, v interface{}) error {
	data, err := read(fsys, name)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// helper function returns true if the runtime engine is
// kubernetes or is container-based.
func isContainerRuntime(pipeline *spec.Pipeline) bool {
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
