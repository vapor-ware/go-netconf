// Copyright (c) 2018, Juniper Networks, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package netconf

import (
	"os/exec"
	"syscall"
)

// TransportJunos maintains the information necessary to communicate with Junos
// via local shell NETCONF interface.
type TransportJunos struct {
	TransportBasicIO
	cmd *exec.Cmd
}

// Close closes an existing local NETCONF session.
func (t *TransportJunos) Close() error {
	if t.cmd != nil {
		t.ReadWriteCloser.Close()
		if err := t.cmd.Process.Kill(); err != nil {
			return err
		}

		if err := t.cmd.Wait(); err != nil {
			return err
		}
	}
	return nil
}

// Open creates a new local NETCONF session.
func (t *TransportJunos) Open() error {
	var err error

	t.cmd = exec.Command("xml-mode", "netconf", "need-trailer")

	t.cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	w, err := t.cmd.StdinPipe()
	if err != nil {
		return err
	}

	r, err := t.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	t.ReadWriteCloser = NewReadWriteCloser(r, w)
	return t.cmd.Start()
}

// DialJunos creates a new NETCONF session via Junos local shell
// NETCONF interface (xml-mode netconf need-trailer).
func DialJunos() (*Session, error) {
	var t TransportJunos
	err := t.Open()
	if err != nil {
		return nil, err
	}
	return NewSession(&t), nil
}
