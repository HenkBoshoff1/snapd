// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package osutil

import (
	"io"
	"os/exec"
	"os/user"
	"syscall"
	"time"
)

func MockUserLookup(mock func(name string) (*user.User, error)) func() {
	realUserLookup := userLookup
	userLookup = mock

	return func() { userLookup = realUserLookup }
}

func MockUserCurrent(mock func() (*user.User, error)) func() {
	realUserCurrent := userCurrent
	userCurrent = mock

	return func() { userCurrent = realUserCurrent }
}

func MockSudoersDotD(mockDir string) func() {
	realSudoersD := sudoersDotD
	sudoersDotD = mockDir

	return func() { sudoersDotD = realSudoersD }
}

func MockMountInfoPath(mockMountInfoPath string) func() {
	realMountInfoPath := mountInfoPath
	mountInfoPath = mockMountInfoPath

	return func() { mountInfoPath = realMountInfoPath }
}

func MockSyscallKill(f func(int, syscall.Signal) error) func() {
	oldSyscallKill := syscallKill
	syscallKill = f
	return func() {
		syscallKill = oldSyscallKill
	}
}

func MockSyscallGetpgid(f func(int) (int, error)) func() {
	oldSyscallGetpgid := syscallGetpgid
	syscallGetpgid = f
	return func() {
		syscallGetpgid = oldSyscallGetpgid
	}
}

func MockCmdWaitTimeout(timeout time.Duration) func() {
	oldCmdWaitTimeout := cmdWaitTimeout
	cmdWaitTimeout = timeout
	return func() {
		cmdWaitTimeout = oldCmdWaitTimeout
	}
}

var KillProcessGroup = killProcessGroup

func WaitingReaderGuts(r io.Reader) (io.Reader, *exec.Cmd) {
	wr := r.(*waitingReader)
	return wr.reader, wr.cmd
}
