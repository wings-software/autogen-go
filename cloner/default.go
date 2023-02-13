// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cloner

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

// New returns a cloner that uses exec.
func New(depth int, out io.Writer) Cloner {
	return &execer{
		depth:  depth,
		stdout: out,
	}
}

// default cloner using the built-in Git client.
type execer struct {
	depth  int
	stdout io.Writer
}

// Clone the repository using the built-in Git client.
func (c *execer) Clone(ctx context.Context, params Params) error {
	// TODO need to set username/ password
	cmd := exec.Command("git", "clone", fmt.Sprintf("--depth=%d", c.depth), params.Repo, params.Dir)
	cmd.Stdout = c.stdout
	cmd.Stderr = c.stdout
	return cmd.Run()
}
