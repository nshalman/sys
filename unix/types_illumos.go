// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

/*
Input to cgo -godefs.  See README.md
*/

package unix

/*
// Many illumos distributions ship a 3rd party tun/tap driver
// from https://github.com/kaizawa/tuntap
// It supports a pair of IOCTLs defined at
// https://github.com/kaizawa/tuntap/blob/master/if_tun.h#L91-L93
#define TUNNEWPPA	(('T'<<16) | 0x0001)
#define TUNSETPPA	(('T'<<16) | 0x0002)
*/
import "C"

const(
	TUNNEWPPA = C.TUNNEWPPA
	TUNSETPPA = C.TUNSETPPA
)
