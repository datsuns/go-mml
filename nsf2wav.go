package main

//#cgo LDFLAGS: ./lib/libGme.a -lstdc++
// int nfs2wav(char* src, char* dest);
import (
	"C"
)

func nsf2wav(src, dest string) int {
	csrc := C.CString(src)
	cdest := C.CString(dest)
	return int(C.nfs2wav(csrc, cdest))
}
