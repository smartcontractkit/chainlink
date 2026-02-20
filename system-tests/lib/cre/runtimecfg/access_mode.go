package runtimecfg

import (
	"fmt"
	"os"
	"strings"
)

const (
	EnvRemoteAccessMode = "CRE_REMOTE_ACCESS_MODE"
	EnvEC2HostIP        = "CRE_EC2_HOST_IP"

	RemoteAccessModeSSM    = "ssm"
	RemoteAccessModeDirect = "direct"
)

func RemoteAccessMode() string {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv(EnvRemoteAccessMode)))
	if mode == "" {
		return RemoteAccessModeSSM
	}
	if mode == RemoteAccessModeDirect {
		return mode
	}
	return RemoteAccessModeSSM
}

func IsDirectMode() bool {
	return RemoteAccessMode() == RemoteAccessModeDirect
}

func DirectHostIP() (string, error) {
	hostIP := strings.TrimSpace(os.Getenv(EnvEC2HostIP))
	if hostIP == "" {
		return "", fmt.Errorf("%s must be set when %s=%s", EnvEC2HostIP, EnvRemoteAccessMode, RemoteAccessModeDirect)
	}
	return hostIP, nil
}
