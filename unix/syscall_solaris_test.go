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
	name := fobj.GetName()
	if path != name {
		t.Errorf(`Can't get name back out: "%s" "%s"`, path, name)
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
