/*
 * GML - Go QML
 *
 * The MIT License (MIT)
 *
 * Copyright (c) 2019 Roland Singer <roland.singer[at]desertbit.com>
 * Copyright (c) 2019 Sebastian Borchers <sebastian[at]desertbit.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package gml

// #include <gml.h>
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type bytes struct {
	freed bool
	ptr   C.gml_bytes
}

func newBytes() (b *bytes) {
	b = &bytes{
		ptr: C.gml_bytes_new(),
	}

	// This should never happen. Signalizes a fatal error.
	if b.ptr == nil {
		panic(fmt.Errorf("failed to create gml bytes: C pointer is nil"))
	}

	runtime.SetFinalizer(b, freeBytes) // Always free the C value.
	return
}

func freeBytes(b *bytes) {
	if b.freed {
		return
	}
	b.freed = true
	C.gml_bytes_free(b.ptr)
}

func (b *bytes) Free() {
	freeBytes(b)
}

func (b *bytes) Bytes() []byte {
	var size C.int
	buf := C.gml_bytes_get(b.ptr, &size)
	if size <= 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(buf), size)
}