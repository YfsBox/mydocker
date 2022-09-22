package image

import (
	"testing"
)

func TestDownLoad(t *testing.T) {
	DownloadImageIfNeed("ubuntu:latest")	
}