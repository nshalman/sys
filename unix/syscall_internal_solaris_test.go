// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build solaris
// +build solaris

package unix

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

func (e *EventPort) checkInternals(t *testing.T, fds, paths, cookies, pending int) {
	t.Helper()
	p, err := e.Pending()
	if err != nil {
		t.Fatalf("failed to query how many events are pending")
	}
	format := "| fds: %d | paths: %d | cookies: %d | pending: %d |"
	expected := fmt.Sprintf(format, fds, paths, cookies, pending)
	got := fmt.Sprintf(format, len(e.fds), len(e.paths), len(e.cookies), p)
	if expected != got {
		t.Errorf("Internal state mismatch\nfound:    %s\nexpected: %s", got, expected)
	}
}

// Regression test for DissociatePath returning ENOENT
// This test is intended to create a linear worst
// case scenario of events being associated and
// fired but not consumed before additional
// calls to dissociate and associate happen
// This needs to be an internal test so that
// we can validate the state of the private maps
func TestEventPortDissociateAlreadyGone(t *testing.T) {
	port, err := NewEventPort()
	if err != nil {
		t.Fatalf("failed to create an EventPort")
	}
	defer port.Close()
	tmpfile, err := ioutil.TempFile("", "eventport")
	if err != nil {
		t.Fatalf("unable to create a tempfile: %v", err)
	}
	path := tmpfile.Name()
	stat, err := os.Stat(path)
	if err != nil {
		t.Fatalf("unexpected failure to Stat file: %v", err)
	}
	err = port.AssociatePath(path, stat, FILE_MODIFIED, "cookie1")
	if err != nil {
		t.Fatalf("unexpected failure associating file: %v", err)
	}
	// We should have 1 path registered and 1 cookie in the jar
	port.checkInternals(t, 0, 1, 1, 0)
	// The path is associated, let's delete it.
	err = os.Remove(path)
	if err != nil {
		t.Fatalf("unexpected failure deleting file: %v", err)
	}
	// The file has been deleted, some sort of pending event is probably
	// queued in the kernel waiting for us to get it AND the kernel is
	// no longer watching for events on it. BUT... Because we haven't
	// consumed the event, this API thinks it's still watched:
	watched := port.PathIsWatched(path)
	if !watched {
		t.Errorf("unexpected result from PathIsWatched")
	}
	// Ok, let's dissociate the file even before reading the event.
	// Oh, ENOENT. I guess it's not associated any more
	err = port.DissociatePath(path)
	if err != ENOENT {
		t.Errorf("unexpected result dissociating a seemingly associated path we know isn't: %v", err)
	}
	// As established by the return value above, this should clearly be false now:
	// Narrator voice: in the first version of this API, it wasn't.
	watched = port.PathIsWatched(path)
	if watched {
		t.Errorf("definitively unwatched file still in the map")
	}
	// We should have nothing registered, but 1 pending event corresponding
	// to the cookie in the jar
	port.checkInternals(t, 0, 0, 1, 1)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatalf("creating test file failed: %s", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("unexpected failure closing file: %v", err)
	}
	stat, err = os.Stat(path)
	c := "cookie2" // c is for cookie, that's good enough for me
	err = port.AssociatePath(path, stat, FILE_MODIFIED, c)
	if err != nil {
		t.Errorf("unexpected failure associating file: %v", err)
	}
	// We should have 1 registered path and its cookie
	// as well as a second cookie corresponding to the pending event
	port.checkInternals(t, 0, 1, 2, 1)

	// Fire another event
	err = os.Remove(path)
	if err != nil {
		t.Fatalf("unexpected failure deleting file: %v", err)
	}
	port.checkInternals(t, 0, 1, 2, 2)
	err = port.DissociatePath(path)
	if err != ENOENT {
		t.Errorf("unexpected result dissociating a seemingly associated path we know isn't: %v", err)
	}
	// By dissociating this path after deletion we ensure that the paths map is now empty
	// If we're not careful we could trigger a nil pointer exception
	port.checkInternals(t, 0, 0, 2, 2)

	f, err = os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		t.Fatalf("creating test file failed: %s", err)
	}
	err = f.Close()
	if err != nil {
		t.Fatalf("unexpected failure closing file: %v", err)
	}
	stat, err = os.Stat(path)
	if err != nil {
		t.Fatalf("unexpected failure to Stat file: %v", err)
	}
	// Put a seemingly duplicate cookie in the jar to see if we can trigger an incorrect removal from the paths map
	err = port.AssociatePath(path, stat, FILE_MODIFIED, c)
	if err != nil {
		t.Errorf("unexpected failure associating file: %v", err)
	}
	port.checkInternals(t, 0, 1, 3, 2)

	// run the garbage collector so that if me messed up it should be painfully clear
	runtime.GC()

	// Before the fix, this would cause a nil pointer exception
	e, err := port.GetOne(nil)
	if err != nil {
		t.Errorf("failed to get an event: %v", err)
	}
	port.checkInternals(t, 0, 1, 2, 1)
	if e.Cookie != "cookie1" {
		t.Errorf(`expected "cookie1", got "%v"`, e.Cookie)
	}
	// Make sure that a cookie of the same value doesn't cause removal from the paths map incorrectly
	e, err = port.GetOne(nil)
	if err != nil {
		t.Errorf("failed to get an event: %v", err)
	}
	port.checkInternals(t, 0, 1, 1, 0)
	if e.Cookie != "cookie2" {
		t.Errorf(`expected "cookie2", got "%v"`, e.Cookie)
	}

	err = os.Remove(path)
	if err != nil {
		t.Fatalf("unexpected failure deleting file: %v", err)
	}
	// Event has fired, but until processed it should still be in the map
	port.checkInternals(t, 0, 1, 1, 1)
	e, err = port.GetOne(nil)
	if err != nil {
		t.Errorf("failed to get an event: %v", err)
	}
	if e.Cookie != "cookie2" {
		t.Errorf(`expected "cookie2", got "%v"`, e.Cookie)
	}
	// The maps should be empty and there should be no pending events
	port.checkInternals(t, 0, 0, 0, 0)
}
