package util

import (
	"runtime"
	"strings"
)

// GetPackPath Get route prefix path for current package
func GetPackPath() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	return "/" + parts[len(parts)-2] + "/"
}
