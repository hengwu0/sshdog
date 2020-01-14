// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: High-level file comment.
package pty

import (
	"io"
	"os"
	"os/exec"
)

type Pty struct {
	pty *os.File
	tty *os.File
}

type ptyWindow struct {
	rows uint16
	cols uint16
	xpix uint16
	ypix uint16
}

func OpenPty() (*Pty, error) {
	pty, tty, err := open_pty()
	if err != nil {
		return nil, err
	}
	return &Pty{pty, tty}, nil
}

func (pty *Pty) CloseTTY() {
	pty.tty.Close()
}

// Execute an exec.Cmd attached to a pty
func (pty *Pty) AttachTty(cmd *exec.Cmd) {
	if cmd.Stdout == nil {
		cmd.Stdout = pty.tty
	}
	if cmd.Stderr == nil {
		cmd.Stderr = pty.tty
	}
	if cmd.Stdin == nil {
		cmd.Stdin = pty.tty
	}
	attach_tty(pty.tty, cmd)
}

// Close the devices
func (pty *Pty) Close() {
	pty.tty.Close()
	pty.pty.Close()
}

// Resize the pty
func (pty *Pty) Resize(rows, cols, xpix, ypix uint16) error {
	win := &ptyWindow{rows, cols, xpix, ypix}
	return resize_pty(pty.pty, win)
}

// Attach to IO
func (pty *Pty) AttachIO(w io.Writer, r io.Reader) {
	//teardown session
	go io.Copy(pty.pty, r)
	go io.Copy(w, pty.pty)
}
