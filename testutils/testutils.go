package testutils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func FilePath(relative string, up ...int) string {
	offset := 1
	if len(up) > 0 && up[0] > 0 {
		offset = up[0]
	}
	_, file, _, ok := runtime.Caller(offset)
	if !ok {
		panic("failed to get file path")
	}
	if filepath.IsAbs(relative) {
		panic(fmt.Sprintf("%s is an absolute file path, must be relative to the current go source file", relative))
	}
	abs := filepath.Join(filepath.Dir(file), relative)
	return abs
}
