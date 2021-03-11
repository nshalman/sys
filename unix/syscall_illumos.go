// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// illumos system calls not present on Solaris.

//go:build amd64 && illumos
// +build amd64,illumos

package unix

import "unsafe"

func bytes2iovec(bs [][]byte) []Iovec {
	iovecs := make([]Iovec, len(bs))
	for i, b := range bs {
		iovecs[i].SetLen(len(b))
		if len(b) > 0 {
			// somehow Iovec.Base on illumos is (*int8), not (*byte)
			iovecs[i].Base = (*int8)(unsafe.Pointer(&b[0]))
		} else {
			iovecs[i].Base = (*int8)(unsafe.Pointer(&_zero))
		}
	}
	return iovecs
}

//sys	readv(fd int, iovs []Iovec) (n int, err error)

func Readv(fd int, iovs [][]byte) (n int, err error) {
	iovecs := bytes2iovec(iovs)
	n, err = readv(fd, iovecs)
	return n, err
}

//sys	preadv(fd int, iovs []Iovec, off int64) (n int, err error)

func Preadv(fd int, iovs [][]byte, off int64) (n int, err error) {
	iovecs := bytes2iovec(iovs)
	n, err = preadv(fd, iovecs, off)
	return n, err
}

//sys	writev(fd int, iovs []Iovec) (n int, err error)

func Writev(fd int, iovs [][]byte) (n int, err error) {
	iovecs := bytes2iovec(iovs)
	n, err = writev(fd, iovecs)
	return n, err
}

//sys	pwritev(fd int, iovs []Iovec, off int64) (n int, err error)

func Pwritev(fd int, iovs [][]byte, off int64) (n int, err error) {
	iovecs := bytes2iovec(iovs)
	n, err = pwritev(fd, iovecs, off)
	return n, err
}

//sys	accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error) = libsocket.accept4

func Accept4(fd int, flags int) (nfd int, sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	var len _Socklen = SizeofSockaddrAny
	nfd, err = accept4(fd, &rsa, &len, flags)
	if err != nil {
		return
	}
	if len > SizeofSockaddrAny {
		panic("RawSockaddrAny too small")
	}
	sa, err = anyToSockaddr(fd, &rsa)
	if err != nil {
		Close(nfd)
		nfd = 0
	}
	return
}

//go:cgo_import_dynamic libc_putmsg putmsg "libc.so"
//go:cgo_import_dynamic libc_getmsg getmsg "libc.so"
//go:linkname f_putmsg libc_putmsg
//go:linkname f_getmsg libc_getmsg

var (
	f_putmsg uintptr
	f_getmsg uintptr
)

func Ioctl(fd int, req uint, arg uintptr) (r int, err error) {
	r0, _, e0 := sysvicall6(uintptr(unsafe.Pointer(&procioctl)), 3,
		uintptr(fd), uintptr(req), arg, 0, 0, 0)

	r = int(r0)
	if e0 != 0 {
		err = e0
	}

	return
}

func Putmsg(fd int, ctlptr uintptr, dataptr uintptr, flags int) (
	r int, err error) {
	r0, _, e0 := sysvicall6(uintptr(unsafe.Pointer(&f_putmsg)), 4,
		uintptr(fd), ctlptr, dataptr, uintptr(flags), 0, 0)

	r = int(r0)
	if e0 != 0 {
		err = e0
	}

	return
}

func Getmsg(fd int, ctlptr uintptr, dataptr uintptr, flagsp uintptr) (
	r int, err error) {
	r0, _, e0 := sysvicall6(uintptr(unsafe.Pointer(&f_getmsg)), 4,
		uintptr(fd), ctlptr, dataptr, flagsp, 0, 0)

	r = int(r0)
	if e0 != 0 {
		err = e0
	}

	return
}
