package fmi2

import (
	"runtime"
	"strings"
)

type Machine struct {
	Architecture  string
	Platform      string
	LibrarySuffix string
}

func CurrentMachine() Machine {

	platform := strings.ToLower(runtime.GOOS)
	machine := strings.ToLower(runtime.GOARCH)

	intSize := 32 << (^uint(0) >> 63) // 32 or 64

	architecture := ""

	suffix := func() string {
		switch platform {
		case "windows":
			return "dll"
		case "darwin":
			return "dylib"
		default:
			return "so"
		}
	}()

	switch machine {
	case "aarch64", "arm64":
		platform += "64"
		architecture = "aarch64"
	case "amd64", "i386", "i686", "x86", "x86_64", "x86pc":
		switch intSize {
		case 32:
			platform += "32"
			architecture = "x86"
		case 64:
			platform += "64"
			architecture = "x86_64"
		}
	}

	return Machine{
		Architecture:  architecture,
		Platform:      platform,
		LibrarySuffix: suffix,
	}
}
