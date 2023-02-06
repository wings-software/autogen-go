// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chroot

import (
	"bytes"
	"testing"
)

func TestGlob(t *testing.T) {
	fs, err := New("testdata/")
	if err != nil {
		t.Error(err)
		return
	}

	matches, err := fs.Glob("/*.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if got, want := len(matches), 2; got != want {
		t.Errorf("Expect %v glob matches, got %v", want, got)
		return
	}

	// verify the base path is stripped from the
	// file names returned in the glob match results.
	if got, want := matches[0], "/en.txt"; got != want {
		t.Errorf("Expect match %v, got %v", want, got)
		return
	}
}

func TestReadFile(t *testing.T) {
	fs, err := New("testdata/")
	if err != nil {
		t.Error(err)
		return
	}

	data, err := fs.ReadFile("/en.txt")
	if err != nil {
		t.Error(err)
		return
	}

	if bytes.Compare([]byte("hello world"), data) != 0 {
		t.Errorf("Expect read file returns file contents")
	}

	if _, err := fs.ReadFile("/es.txt"); err == nil {
		t.Errorf("Expect error when file does not exist")
	}
}

func TestStat(t *testing.T) {
	fs, err := New("testdata/")
	if err != nil {
		t.Error(err)
		return
	}

	// ensure stat with absolute path
	_, err = fs.Stat("/en.txt")
	if err != nil {
		t.Error(err)
		return
	}

	// ensure stat with relative path
	_, err = fs.Stat("fr.txt")
	if err != nil {
		t.Error(err)
		return
	}
}
