// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os/exec"
	"strings"
)

type RunError struct {
	Cmd    string
	Args   []string
	Err    error
	Stderr []byte
}

func (e *RunError) Error() string {
	text := e.Cmd
	if len(e.Args) > 0 {
		text += " " + strings.Join(e.Args, " ")
	}
	text += ": " + e.Err.Error()

	if stderr := bytes.TrimRight(e.Stderr, "\n"); len(stderr) > 0 {
		text += ":\n\t" + strings.Replace(string(stderr), "\n", "\n\t", -1)
	}

	return text
}

func Run(dir, name string, args ...string) ([]byte, error) {
	c := exec.Command(name, args...)
	c.Dir = dir

	var stdout, stderr bytes.Buffer
	c.Stdout, c.Stderr = &stdout, &stderr

	err := c.Run()
	if err != nil {
		err = &RunError{
			Cmd:    name,
			Args:   args,
			Stderr: stderr.Bytes(),
			Err:    err,
		}
	}

	return stdout.Bytes(), err
}
