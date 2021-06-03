// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build solaris
// +build solaris

package unix_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"testing"

	"golang.org/x/sys/unix"
)

func TestStatvfs(t *testing.T) {
	if err := unix.Statvfs("", nil); err == nil {
		t.Fatal(`Statvfs("") expected failure`)
	}

	statvfs := unix.Statvfs_t{}
	if err := unix.Statvfs("/", &statvfs); err != nil {
		t.Errorf(`Statvfs("/") failed: %v`, err)
	}

	if t.Failed() {
		mount, err := exec.Command("mount").CombinedOutput()
		if err != nil {
			t.Logf("mount: %v\n%s", err, mount)
		} else {
			t.Logf("mount: %s", mount)
		}
	}
}

func TestSysconf(t *testing.T) {
	n, err := unix.Sysconf(3 /* SC_CLK_TCK */)
	if err != nil {
		t.Fatalf("Sysconf: %v", err)
	}
	t.Logf("Sysconf(SC_CLK_TCK) = %d", n)
}

// Event Ports

func TestBasicEventPort(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "eventport")
	if err != nil {
		t.Errorf("unable to create a tempfile: %v", err)
	}
	path := tmpfile.Name()
	defer os.Remove(path)

	stat, err := os.Stat(path)
	if err != nil {
		t.Errorf("Failed to stat %s: %v", path, err)
	}
	port, err := unix.NewEventPort()
	if err != nil {
		t.Errorf("NewEventPort failed: %v", err)
	}
	defer port.Close()
	var cookie unix.EventPortUserCookie = stat.Mode()
	err = port.AssociatePath(path, stat, unix.FILE_MODIFIED, &cookie)
	if err != nil {
		t.Errorf("AssociatePath failed: %v", err)
	}
	if !port.PathIsWatched(path) {
		t.Errorf("PathIsWatched unexpectedly returned false")
	}
	err = port.DissociatePath(path)
	if err != nil {
		t.Errorf("DissociatePath failed: %v", err)
	}
	err = port.AssociatePath(path, stat, unix.FILE_MODIFIED, &cookie)
	if err != nil {
		t.Errorf("AssociatePath failed: %v", err)
	}
	bs := []byte{42}
	tmpfile.Write(bs)
	timeout := new(unix.Timespec)
	timeout.Sec = 1
	pevent, err := port.Get(timeout)
	if err == unix.ETIME {
		t.Errorf("PortGet timed out: %v", err)
	}
	if err != nil {
		t.Errorf("PortGet failed: %v", err)
	}
	if pevent.Path != path {
		t.Errorf("Path mismatch: %v != %v", pevent.Path, path)
	}
	err = port.AssociatePath(path, stat, unix.FILE_MODIFIED, &cookie)
	if err != nil {
		t.Errorf("AssociatePath failed: %v", err)
	}
	err = port.AssociatePath(path, stat, unix.FILE_MODIFIED, &cookie)
	if err == nil {
		t.Errorf("Unexpected success associating already associated path")
	}
}

func TestEventPortFds(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	stat, err := os.Stat(path)
	fmode := stat.Mode()
	port, err := unix.NewEventPort()
	if err != nil {
		t.Errorf("NewEventPort failed: %v", err)
	}
	defer port.Close()
	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("unable to create a pipe: %v", err)
	}
	defer w.Close()
	defer r.Close()
	fd := r.Fd()

	var cookie unix.EventPortUserCookie = fmode
	port.AssociateFd(fd, unix.POLLIN, &cookie)
	if !port.FdIsWatched(fd) {
		t.Errorf("FdIsWatched unexpectedly returned false")
	}
	err = port.DissociateFd(fd)
	err = port.AssociateFd(fd, unix.POLLIN, &cookie)
	bs := []byte{42}
	w.Write(bs)
	timeout := new(unix.Timespec)
	timeout.Sec = 1
	pevent, err := port.Get(timeout)
	if err == unix.ETIME {
		t.Errorf("PortGet timed out: %v", err)
	}
	if err != nil {
		t.Errorf("PortGet failed: %v", err)
	}
	if pevent.Fd != fd {
		t.Errorf("Fd mismatch: %v != %v", pevent.Fd, fd)
	}
	var c = pevent.Cookie
	if c == nil {
		t.Errorf("Cookie missing: %v != %v", &cookie, c)
		return
	}
	if *c != cookie {
		t.Errorf("Cookie mismatch: %v != %v", cookie, *c)
	}
	port.AssociateFd(fd, unix.POLLIN, &cookie)
	err = port.AssociateFd(fd, unix.POLLIN, &cookie)
	if err == nil {
		t.Errorf("unexpected success associating already associated fd")
	}
}

func TestEventPortErrors(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "eventport")
	if err != nil {
		t.Errorf("unable to create a tempfile: %v", err)
	}
	path := tmpfile.Name()
	stat, _ := os.Stat(path)
	os.Remove(path)
	port, _ := unix.NewEventPort()
	err = port.AssociatePath(path, stat, unix.FILE_MODIFIED, nil)
	if err == nil {
		t.Errorf("unexpected success associating nonexistant file")
	}
	err = port.DissociatePath(path)
	if err == nil {
		t.Errorf("unexpected success dissociating unassociated path")
	}
	timeout := new(unix.Timespec)
	timeout.Nsec = 1
	_, err = port.Get(timeout)
	if err != unix.ETIME {
		t.Errorf("unexpected lack of timeout")
	}
	err = port.DissociateFd(uintptr(0))
	if err == nil {
		t.Errorf("unexpected success dissociating unassociated fd")
	}
}
