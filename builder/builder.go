// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package builder builds a pipeline configuration.
package builder

import (
	"io/fs"

	spec "github.com/drone/spec/dist/go"
	"github.com/ghodss/yaml"
	"github.com/wings-software/autogen-go/utils"
)

// Builder builds a pipeline configuration.
type Builder struct {
	vendor Vendor
}

// New creates a new pipeline builder.
func New(vendor string) *Builder {
	return &Builder{
		vendor: NewVendor(vendor),
	}
}

// Build the pipeline configuration.
func (b *Builder) Build(fsys fs.FS) ([]byte, error) {
	stageci := new(spec.StageCI)
	// removing it as its not required.
	// stageci.Platform = new(spec.Platform)
	// stageci.Platform.Os = spec.OSLinux
	// stageci.Platform.Arch = spec.ArchAmd64

	stage := new(spec.Stage)
	stage.Name = "build"
	stage.Type = "ci"
	stage.Spec = stageci

	pipeline := new(spec.Pipeline)
	pipeline.Name = "sample pipeline"
	pipeline.Stages = append(pipeline.Stages, stage)
	for _, rule := range b.vendor.GetRules() {
		if err := rule(fsys, pipeline); err == utils.SkipAll {
			break
		}

		// we purposefully ignore errors here.
		// an error in an individual rule should
		// never prevent yaml generation.
	}

	return yaml.Marshal(pipeline)
}
