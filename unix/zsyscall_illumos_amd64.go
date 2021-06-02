// go run mksyscall_solaris.go -illumos -tags illumos,amd64 syscall_illumos.go
// Code generated by the command above; see README.md. DO NOT EDIT.

//go:build illumos && amd64
// +build illumos,amd64

package unix

import (
	"unsafe"
)

//go:cgo_import_dynamic libc_readv readv "libc.so"
//go:cgo_import_dynamic libc_preadv preadv "libc.so"
//go:cgo_import_dynamic libc_writev writev "libc.so"
//go:cgo_import_dynamic libc_pwritev pwritev "libc.so"
//go:cgo_import_dynamic libc_accept4 accept4 "libsocket.so"
//go:cgo_import_dynamic libc_putmsg putmsg "libc.so"
//go:cgo_import_dynamic libc_getmsg getmsg "libc.so"
//go:cgo_import_dynamic libc_port_create port_create "libc.so"
//go:cgo_import_dynamic libc_port_associate port_associate "libc.so"
//go:cgo_import_dynamic libc_port_dissociate port_dissociate "libc.so"
//go:cgo_import_dynamic libc_port_get port_get "libc.so"

//go:linkname procreadv libc_readv
//go:linkname procpreadv libc_preadv
//go:linkname procwritev libc_writev
//go:linkname procpwritev libc_pwritev
//go:linkname procaccept4 libc_accept4
//go:linkname procputmsg libc_putmsg
//go:linkname procgetmsg libc_getmsg
//go:linkname procport_create libc_port_create
//go:linkname procport_associate libc_port_associate
//go:linkname procport_dissociate libc_port_dissociate
//go:linkname procport_get libc_port_get

var (
	procreadv,
	procpreadv,
	procwritev,
	procpwritev,
	procaccept4,
	procputmsg,
	procgetmsg,
	procport_create,
	procport_associate,
	procport_dissociate,
	procport_get syscallFunc
)

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func readv(fd int, iovs []Iovec) (n int, err error) {
	var _p0 *Iovec
	if len(iovs) > 0 {
		_p0 = &iovs[0]
	}
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procreadv)), 3, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(len(iovs)), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func preadv(fd int, iovs []Iovec, off int64) (n int, err error) {
	var _p0 *Iovec
	if len(iovs) > 0 {
		_p0 = &iovs[0]
	}
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procpreadv)), 4, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(len(iovs)), uintptr(off), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func writev(fd int, iovs []Iovec) (n int, err error) {
	var _p0 *Iovec
	if len(iovs) > 0 {
		_p0 = &iovs[0]
	}
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procwritev)), 3, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(len(iovs)), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func pwritev(fd int, iovs []Iovec, off int64) (n int, err error) {
	var _p0 *Iovec
	if len(iovs) > 0 {
		_p0 = &iovs[0]
	}
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procpwritev)), 4, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(len(iovs)), uintptr(off), 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func accept4(s int, rsa *RawSockaddrAny, addrlen *_Socklen, flags int) (fd int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procaccept4)), 4, uintptr(s), uintptr(unsafe.Pointer(rsa)), uintptr(unsafe.Pointer(addrlen)), uintptr(flags), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func putmsg(fd int, clptr *strbuf, dataptr *strbuf, flags int) (err error) {
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procputmsg)), 4, uintptr(fd), uintptr(unsafe.Pointer(clptr)), uintptr(unsafe.Pointer(dataptr)), uintptr(flags), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func getmsg(fd int, clptr *strbuf, dataptr *strbuf, flags *int) (err error) {
	_, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procgetmsg)), 4, uintptr(fd), uintptr(unsafe.Pointer(clptr)), uintptr(unsafe.Pointer(dataptr)), uintptr(unsafe.Pointer(flags)), 0, 0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func port_create() (n int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procport_create)), 0, 0, 0, 0, 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func port_associate(port int, source int, object uintptr, events int, user unsafe.Pointer) (n int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procport_associate)), 5, uintptr(port), uintptr(source), uintptr(object), uintptr(events), uintptr(user), 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func port_dissociate(port int, source int, object uintptr) (n int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procport_dissociate)), 3, uintptr(port), uintptr(source), uintptr(object), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}

// THIS FILE IS GENERATED BY THE COMMAND AT THE TOP; DO NOT EDIT

func port_get(port int, pe *PortEvent, timeout *Timespec) (n int, err error) {
	r0, _, e1 := sysvicall6(uintptr(unsafe.Pointer(&procport_get)), 3, uintptr(port), uintptr(unsafe.Pointer(pe)), uintptr(unsafe.Pointer(timeout)), 0, 0, 0)
	n = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}
