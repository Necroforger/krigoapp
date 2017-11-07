package window

// #include "window.h"
import "C"
import (
	"bytes"
	"encoding/binary"
	"unicode/utf16"
	"unsafe"
)

// HWND holds window data
type HWND struct {
	hwnd C.HWND
}

// Title retrieves the title of a window
func (h *HWND) Title() string {
	title := C.windowText(h.hwnd)
	gostr := WCHARToGoString(title)
	C.free(unsafe.Pointer(title))
	return gostr
}

// GetForegroundWindow ...
func GetForegroundWindow() *HWND {
	return &HWND{
		hwnd: C.foregroundWindow(),
	}
}

// UTF16ToString ...
func UTF16ToString(s []uint16) string {
	for i, v := range s {
		if v == 0 {
			s = s[0:i]
			break
		}
	}
	return string(utf16.Decode(s))
}

// WCHARToGoString converts a wchar to a go string
func WCHARToGoString(str *C.WCHAR) string {
	n := int(C.wcharlen(str))
	data := C.GoBytes(unsafe.Pointer(str), C.int(n*16))

	encoded := make([]uint16, n)
	buf := bytes.NewBuffer(data)
	for idx := 0; idx < n; idx++ {
		err := binary.Read(buf, binary.LittleEndian, &encoded[idx])
		if err != nil {
			break
		}
	}

	gostr := UTF16ToString(encoded)
	return gostr
}
