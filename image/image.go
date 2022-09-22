package image

import (
	"fmt"
	cm "mydocker/common"
	"os"
	"github.com/haolipeng/archiver"
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

	fmt.Printf("the image is %v the imageHexHash is %v",image,imageHexHash)

	if err := saveImageLocal(image,ImageName,imageHexHash); err != nil {
		return fmt.Errorf("saveImageLocal error:%v",err)
	}

	if err := tarImageFiles(imageHexHash); err != nil {
		return fmt.Errorf("tarImageFiles error: %v",err)
	}

	return nil
}


func saveImageLocal(img v1.Image,src string,imgHash string) error {
	
	imageSavePath := fmt.Sprintf("%v/%v",cm.GetTmpPath(),imgHash)
	if err := os.MkdirAll(imageSavePath,757); err != nil {
		return fmt.Errorf("Mkdir imgSavePath error: %v",err)
	} //创建文件夹
	imagePath := imageSavePath + "/package.tar"
	cm.DPrintf("save legacy,src %v,path: %v\n",src,imagePath)
	if err := crane.SaveLegacy(img,src,imagePath); err != nil {
		return fmt.Errorf("SaveLegacy error %v",err)
	}
	cm.DPrintf("Save ok")

	return nil

}

func tarImageFiles(imgHash string) error {
    //构造存储路径
    path := fmt.Sprintf("%v/%v",cm.GetTmpPath(),imgHash)
    tarPath := path + "/package.tar"

    if err := archiver.Untar(tarPath, path); err != nil {
		return fmt.Errorf("Untar %v err %v",tarPath,err)
	}
	return nil
}
