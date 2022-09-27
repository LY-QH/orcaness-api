package util

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/rs/xid"
)

// GetPackPath Get route prefix path for current package
func GetPackPath() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	fmt.Println(parts)
	routerName := parts[len(parts)-1]
	return "/" + routerName[:len(routerName)-7] + "/"
}

func GenId(prefix ...string) string {
	id := xid.New().String()
	if len(prefix) == 0 {
		return id
	}

	return prefix[0] + id
}
