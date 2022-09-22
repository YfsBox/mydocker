package image

import (
	"fmt"
	cm "mydocker/common"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/haolipeng/archiver"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

type Manifest struct {
	Config string	`json:"Config"`
	Tags []string	`json:"RepoTags"`
	Layers []string	`json:"Layers"`
}

func ProcessLayers(ImageHash string) error {

	layerPath := fmt.Sprintf("%v/%v/%v",cm.GetTmpPath(),ImageHash,"layer")
	if err := os.MkdirAll(layerPath,0757); err != nil {
		return fmt.Errorf("mkdir %v error %v",layerPath,err)
	}

	manifestPath := fmt.Sprintf("%v/%v/manifest.json",cm.GetTmpPath(),ImageHash)
	cm.DPrintf("%v",manifestPath)
	manis := []Manifest{}

	content,_ := ioutil.ReadFile(manifestPath)
	cm.DPrintf("begin unmarshal")
	if err := json.Unmarshal(content,&manis); err != nil {
		cm.DPrintf("unmarshal error %v",err)
		return fmt.Errorf("the json parse to mani error:%v",err)
	}
	//将minis中的各layers逐层进行解析

	for _,layer := range manis[0].Layers {
		layerFile := fmt.Sprintf("%v/%v/%v",cm.GetTmpPath(),ImageHash,layer)
		dstPath := fmt.Sprintf("%v/%v",layerPath,layer[:12])

		cm.DPrintf("The layerFile is %v,and the dstPath is %v",layerFile,dstPath)
		if err := os.MkdirAll(dstPath,0757); err != nil { //首先创建目标文件夹,位于layer文件夹之下
			return fmt.Errorf("Mkdir %v error %v",dstPath,err)
		}
		if err := cm.Untar(dstPath,layerFile); err != nil { //这个地方有点问题
			return fmt.Errorf("Untar %v error %v",layerFile,err)
		}
	}

	manifestImagePath := fmt.Sprintf("%v/%v/manifest.json",cm.GetImagePath(),ImageHash)
	/*if err := os.CopyFile(manifestPath,manifestImagePath); err != nil {
		return fmt.Errorf("Copy %v error %v",manifestPath,err)
	}*/
	cm.DPrintf("manifestImagePath: %v",manifestImagePath)

	return nil
}


func DownloadImageIfNeed(ImageName string) (string,error) {
	if ImageName == "" {
		return "",fmt.Errorf("ImageName can't be empty!")
	}
	var image v1.Image
	var err error
	if image,err = crane.Pull(ImageName); err != nil {
		return "",fmt.Errorf("Pull %v err!",ImageName)
	}

	m, err := image.Manifest() //获取镜像的hash值
	imageFullHash := m.Config.Digest.Hex
	imageHexHash := imageFullHash[:12]

	fmt.Printf("the image is %v the imageHexHash is %v",image,imageHexHash)

	if err := saveImageLocal(image,ImageName,imageHexHash); err != nil {
		return "",fmt.Errorf("saveImageLocal error:%v",err)
	}

	if err := tarImageFiles(imageHexHash); err != nil {
		return "",fmt.Errorf("tarImageFiles error: %v",err)
	}

	return imageHexHash,nil
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
