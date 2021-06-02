// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build illumos
// +build illumos

package unix_test

import (
	"golang.org/x/sys/unix"
	"os"
	"runtime"
	"testing"
	"unsafe"
)

func TestLifreqSetName(t *testing.T) {
	var l unix.Lifreq
	err := l.SetName("12345678901234356789012345678901234567890")
	if err == nil {
		t.Fatal(`Lifreq.SetName should reject names that are too long`)
	}
	err = l.SetName("tun0")
	if err != nil {
		t.Errorf(`Lifreq.SetName("tun0") failed: %v`, err)
	}
}

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
	_, err = unix.PortAssociateFileObj(port, fobj, unix.FILE_MODIFIED, unsafe.Pointer(&fmode))
	name := fobj.GetName()
	if path != name {
		t.Errorf(`Can't get name back out: "%s" "%s"`, path, name)
	}
	if err != nil {
		t.Errorf("PortAssociateFileObj failed: %v", err)
	}
}
