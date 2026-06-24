package platform

import (
	"runtime"
)

const (
	OSWindows = "windows"
	OSLinux   = "linux"
	OSMac     = "darwin"
)

const (
	ArchAMD64 = "amd64"
	Arch386   = "386"
	ArchARM64 = "arm64"
	ArchARM   = "arm"
)

type Info struct {
	OS   string
	Arch string
}

func Detect() Info {
	return Info{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

func (info Info) GodotAssetSuffix() string {
	switch info.OS {
	case OSWindows:
		switch info.Arch {
		case ArchAMD64:
			return "win64.exe.zip"
		case Arch386:
			return "win32.exe.zip"
		case ArchARM64:
			return "windows_arm64.exe.zip"
		}
	case OSLinux:
		switch info.Arch {
		case ArchAMD64:
			return "linux.x86_64.zip"
		case Arch386:
			return "linux.x86_32.zip"
		case ArchARM64:
			return "linux.arm64.zip"
		case ArchARM:
			return "linux.arm32.zip"
		}
	case OSMac:
		return "macos.universal.zip"
	}
	return ""
}

func (info Info) GodotMonoAssetSuffix() string {
	switch info.OS {
	case OSWindows:
		switch info.Arch {
		case ArchAMD64:
			return "mono_win64.zip"
		case Arch386:
			return "mono_win32.zip"
		case ArchARM64:
			return "mono_windows_arm64.zip"
		}
	case OSLinux:
		switch info.Arch {
		case ArchAMD64:
			return "mono_linux_x86_64.zip"
		case Arch386:
			return "mono_linux_x86_32.zip"
		case ArchARM64:
			return "mono_linux_arm64.zip"
		case ArchARM:
			return "mono_linux_arm32.zip"
		}
	case OSMac:
		return "mono_macos.universal.zip"
	}
	return ""
}

func (info Info) IsSupported() bool {
	return info.GodotAssetSuffix() != ""
}
