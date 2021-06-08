// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build solaris
// +build solaris

package unix_test

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
	"unsafe"

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

func TestCreateFileObj(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	stat, err := os.Stat(path)
	if err != nil {
		t.Errorf("Failed to stat %s: %v", path, err)
	}
	fobj, err := unix.CreateFileObj(path, stat)
	p := fobj.Path()
	if path != p {
		t.Errorf(`Can't get path back out: "%s" "%s"`, path, p)
	}
}

func TestBasicEventPort(t *testing.T) {
	_, path, _, _ := runtime.Caller(0)
	stat, err := os.Stat(path)
	fmode := stat.Mode()
	if err != nil {
		t.Errorf("Failed to stat %s: %v", path, err)
	}
	port, err := unix.PortCreate()
	if err != nil {
		t.Errorf("PortCreate failed: %d - %v", port, err)
	}
	defer unix.Close(port)
	fobj, err := unix.CreateFileObj(path, stat)
	if err != nil {
		t.Errorf("CreateFileObj failed: %v", err)
	}
	_, err = unix.PortAssociateFileObj(port, fobj, unix.FILE_MODIFIED, (*byte)(unsafe.Pointer(&fmode)))
	if err != nil {
		t.Errorf("PortAssociateFileObj failed: %v", err)
	}
	_, err = unix.PortDissociateFileObj(port, fobj)
	if err != nil {
		t.Errorf("PortDissociateFileObj failed: %v", err)
	}
}

func TestEventPortFds(t *testing.T) {
	port, err := unix.PortCreate()
	if err != nil {
		t.Errorf("PortCreate failed: %d - %v", port, err)
	}
	defer unix.Close(port)
	r, w, err := os.Pipe()
	if err != nil {
		t.Errorf("unable to create a pipe: %v", err)
	}
	defer w.Close()
	defer r.Close()

	unix.PortAssociateFd(port, int(r.Fd()), unix.POLLIN, nil)
	defer unix.PortDissociateFd(port, int(r.Fd()))
	bs := []byte{42}
	w.Write(bs)
	timeout := new(unix.Timespec)
	timeout.Sec = 1
	var pevent unix.PortEvent
	_, err = unix.PortGet(port, &pevent, timeout)
	if err == unix.ETIME {
		t.Errorf("PortGet timed out: %v", err)
	}
	if err != nil {
		t.Errorf("PortGet failed: %v", err)
	}
	fd, err := pevent.Fd()
	if err != nil {
		t.Errorf("Unable to retrieve Fd from PortEvent: %v", err)
	}
	if fd != int(r.Fd()) {
		t.Errorf("Fd mismatch: %v != %v", fd, int(r.Fd()))
	}
}
