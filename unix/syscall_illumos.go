// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// illumos system calls not present on Solaris.

//go:build amd64 && illumos
// +build amd64,illumos

package unix

import (
	"fmt"
	"runtime"
	"unsafe"
)

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

//sys	putmsg(fd int, clptr *strbuf, dataptr *strbuf, flags int) (err error)

func Putmsg(fd int, cl []byte, data []byte, flags int) (err error) {
	var clp, datap *strbuf
	if len(cl) > 0 {
		clp = &strbuf{
			Len: int32(len(cl)),
			Buf: (*int8)(unsafe.Pointer(&cl[0])),
		}
	}
	if len(data) > 0 {
		datap = &strbuf{
			Len: int32(len(data)),
			Buf: (*int8)(unsafe.Pointer(&data[0])),
		}
	}
	return putmsg(fd, clp, datap, flags)
}

//sys	getmsg(fd int, clptr *strbuf, dataptr *strbuf, flags *int) (err error)

func Getmsg(fd int, cl []byte, data []byte) (retCl []byte, retData []byte, flags int, err error) {
	var clp, datap *strbuf
	if len(cl) > 0 {
		clp = &strbuf{
			Maxlen: int32(len(cl)),
			Buf:    (*int8)(unsafe.Pointer(&cl[0])),
		}
	}
	if len(data) > 0 {
		datap = &strbuf{
			Maxlen: int32(len(data)),
			Buf:    (*int8)(unsafe.Pointer(&data[0])),
		}
	}

	if err = getmsg(fd, clp, datap, &flags); err != nil {
		return nil, nil, 0, err
	}

	if len(cl) > 0 {
		retCl = cl[:clp.Len]
	}
	if len(data) > 0 {
		retData = data[:datap.Len]
	}
	return retCl, retData, flags, nil
}

func IoctlPlink(fd, other_fd int) (int, error) {
	muxid, err := ioctlRet(fd, I_PLINK, uintptr(other_fd))
	if err != nil {
		return -1, err
	}

	return muxid, nil
}

func IoctlPunlink(fd, muxid int) error {
	return ioctl(fd, I_PUNLINK, uintptr(muxid))
}

func IoctlGetIPMuxID(fd int, name string) (int, error) {
	var req lifreq
	if len(name) >= len(req.Name) {
		return -1, fmt.Errorf("name cannot be more than %d characters", len(req.Name)-1)
	}
	for i := range name {
		req.Name[i] = int8(name[i])
	}

	// In the ioctl syscall definition, req is a uint, but on Illumos
	// it's an int. This means that the SIOCGLIFMUXID constant is
	// defined as negative, and can't be used inline in the ioctl
	// call. We have to explicitly initialize an int and then cast
	// that to uint.
	reqnum := int(SIOCGLIFMUXID)
	if err := ioctl(fd, uint(reqnum), uintptr(unsafe.Pointer(&req))); err != nil {
		return -1, err
	}

	id := *(*int)(unsafe.Pointer(&req.Lifru[0]))
	return id, nil
}

func IoctlSetIPMuxID(fd int, name string, muxID int) error {
	var req lifreq
	if len(name) >= len(req.Name) {
		return fmt.Errorf("name cannot be more than %d characters", len(req.Name)-1)
	}
	for i := range name {
		req.Name[i] = int8(name[i])
	}
	*(*int)(unsafe.Pointer(&req.Lifru[0])) = muxID

	// In the ioctl syscall definition, req is a uint, but on Illumos
	// it's an int. This means that the SIOCSLIFMUXID constant is
	// defined as negative, and can't be used inline in the ioctl
	// call. We have to explicitly initialize an int and then cast
	// that to uint.
	reqnum := int(SIOCSLIFMUXID)
	err := ioctl(fd, uint(reqnum), uintptr(unsafe.Pointer(&req)))
	runtime.KeepAlive(&req)
	return err
}

func IoctlSetString(fd int, req uint, val string) error {
	bs := make([]byte, len(val)+1)
	copy(bs[:len(bs)-1], val)
	err := ioctl(fd, req, uintptr(unsafe.Pointer(&bs[0])))
	runtime.KeepAlive(&bs[0])
	return err
}

func IoctlTunNewPPA(fd int, ppaNum int) (ppa int, err error) {
	var req strioctl
	req.Cmd = TUNNEWPPA
	req.Len = int32(unsafe.Sizeof(ppaNum))
	req.Dp = (*int8)(unsafe.Pointer(&ppaNum))

	ppa, err = ioctlRet(fd, I_STR, uintptr(unsafe.Pointer(&req)))
	runtime.KeepAlive(&req)
	return ppa, err
}
