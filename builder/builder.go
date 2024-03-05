// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package builder builds a pipeline configuration.
package builder

import (
	"io/fs"

	"github.com/drone/go-convert/convert/harness/downgrader"
	spec "github.com/drone/spec/dist/go"
	"github.com/ghodss/yaml"
	"github.com/wings-software/autogen-go/utils"
)

// Builder builds a pipeline configuration.
type Builder struct {
	vendor  Vendor
	version string
}

// New creates a new pipeline builder.
func New(vendor, version string) *Builder {
	if version == "" {
		version = "v1"
	}
	return &Builder{
		vendor:  NewVendor(vendor),
		version: version,
	}
}

// Build the pipeline configuration.
func (b *Builder) Build(fsys fs.FS) ([]byte, error) {
	stageci := new(spec.StageCI)
	// removing it as its not required.
	// stageci.Platform = new(spec.Platform)
	// stageci.Platform.Os = spec.OSLinux
	// stageci.Platform.Arch = spec.ArchAmd64

	stageci.Cache = new(spec.Cache)
	stageci.Cache.Enabled = true
	stage := new(spec.Stage)
	stage.Name = "build"
	stage.Type = "ci"
	stage.Spec = stageci

	pipeline := new(spec.Pipeline)
	pipeline.Stages = append(pipeline.Stages, stage)
	for _, rule := range b.vendor.GetRules() {
		if err := rule(fsys, pipeline); err == utils.SkipAll {
			break
		}

		// we purposefully ignore errors here.
		// an error in an individual rule should
		// never prevent yaml generation.
	}

	config := new(spec.Config)
	config.Type = "pipeline"
	config.Kind = "pipeline"
	config.Version = 1
	config.Spec = pipeline

	yml, err := yaml.Marshal(config)
	if (b.version == "v1") || (b.version == "default") {
		return yml, err
	} else {
		d := downgrader.New()
		v0Yaml, err := d.Downgrade(yml)
		return v0Yaml, err
	}
}
