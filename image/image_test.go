package image

import (
	cm "mydocker/common"
	"testing"
)

func TestDownLoad(t *testing.T) {
	DownloadImageIfNeed("madduci/docker-linux-cpp")
}

func TestProcessLayers(t *testing.T) {
	cm.InitMyDockerDirs()
	hash, _ := DownloadImageIfNeed("busybox")
	ProcessLayers(hash)

}