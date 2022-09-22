package image

import (
	"fmt"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func DownloadImageIfNeed(ImageName string) error {
	if ImageName == "" {
		return fmt.Errorf("ImageName can't be empty!")
	}
	var image v1.Image
	var err error
	if image,err = crane.Pull(ImageName); err != nil {
		return fmt.Errorf("Pull %v err!",ImageName)
	}

	m, err := image.Manifest() //获取镜像的hash值
	imageFullHash := m.Config.Digest.Hex
	imageHexHash := imageFullHash[:12]

	fmt.Printf("the imageHexHash is %v",imageHexHash)

	return nil
}


func saveImageLocal() error {
	return nil

}

func tarImageFiles() error {
	return nil

}
