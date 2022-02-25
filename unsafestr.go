// Copyright 2021 Intuitive Labs GmbH. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package unsafeconv

import (
	"reflect"
	"unsafe"
)

// reflect.SliceHeader but with unsafe.Pointer instead of uintptr
// (gc-safe to manipulate as opposed to the reflect version)
type sliceHeader struct {
	data unsafe.Pointer
	len  int
	cap  int
}

// reflect.StringHeader but with unsafe.Pointer instead of uintptr
// (gc-safe to manipulate as opposed to the reflect version)
type stringHeader struct {
	data unsafe.Pointer
	len  int
}

func init() {
	// sanity checks, in case future go versions change
	// the string or []byte header format
	// see https://groups.google.com/g/golang-nuts/c/Zsfk-VMd_fU/m/WXPjfZwPBAAJ
	var sliceHdr sliceHeader
	var reflectSliceHdr reflect.SliceHeader
	if unsafe.Sizeof(sliceHdr) != unsafe.Sizeof(reflectSliceHdr) {
		panic("slice internal format has changed: different size")
	}
	if unsafe.Offsetof(sliceHdr.data) !=
		unsafe.Offsetof(reflectSliceHdr.Data) {
		panic("slice internal format has changed: different data offset")
	}
	if unsafe.Offsetof(sliceHdr.len) !=
		unsafe.Offsetof(reflectSliceHdr.Len) {
		panic("slice internal format has changed: different len offset")
	}
	if unsafe.Offsetof(sliceHdr.cap) !=
		unsafe.Offsetof(reflectSliceHdr.Cap) {
		panic("slice internal format has changed: different cap offset")
	}

	var strHdr stringHeader
	var reflectStrHdr reflect.StringHeader
	if unsafe.Sizeof(strHdr) != unsafe.Sizeof(reflectStrHdr) {
		panic("string internal format has changed")
	}
	if unsafe.Offsetof(strHdr.data) != unsafe.Offsetof(reflectStrHdr.Data) {
		panic("string internal format has changed: different data offset")
	}
	if unsafe.Offsetof(strHdr.len) != unsafe.Offsetof(reflectStrHdr.Len) {
		panic("string internal format has changed: different len offset")
	}

	// check direct (*string)(*[]byte)
	if unsafe.Offsetof(reflectSliceHdr.Data) !=
		unsafe.Offsetof(reflectStrHdr.Data) {
		panic("string internal data pointer different from slice")
	}
	if unsafe.Offsetof(reflectSliceHdr.Len) !=
		unsafe.Offsetof(reflectStrHdr.Len) {
		panic("string internal data pointer different from slice")
	}
}

// Str converts a byte slice to a string without making any copy or allocations.// The content of the underlying byte slice _must_ not be changed.
func Str(b []byte) (s string) {
	s = *(*string)(unsafe.Pointer(&b))
	return
}

// Bytes converts a string to a byte slice without any copy or allocations.
// The content of the resulting byte slice _must_ not be changed.
func Bytes(s string) []byte {
	const MaxInt32 = 1<<31 - 1
	// convert s data pointer to huge max byte array pointer ([MaxInt32]byte)
	// and then take a slice of it.
	// len(s) & MaxInt32 makes sure the array bounds checks are optimised
	// away.
	// see
	// https://stackoverflow.com/a/69231355 and
	// https://groups.google.com/g/golang-nuts/c/Zsfk-VMd_fU/m/O1ru4fO-BgAJ

	return (*[MaxInt32]byte)((*stringHeader)(
		unsafe.Pointer(&s)).data)[: len(s)&MaxInt32 : len(s)&MaxInt32]

	/* alternative:
	bhdr := (*sliceHeader)(unsafe.Pointer(&b))
	shdr := (*strHeader)(unsafe.Pointer(&s))
	bhdr.data = shdr.data
	bhdr.len = len(s)
	bhdr.cap = len(s)
	*/
	/* alternative go 1.17+:
	unsafe.Slice((*byte)( (*stringHeader)(unsafe.Pointer(&s)).data), len(s))
	*/
}
