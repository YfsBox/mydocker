package image

import (
	"testing"
)

func TestDownLoad(t *testing.T) {
	DownloadImageIfNeed("madduci/docker-linux-cpp")	
}