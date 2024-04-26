package dockerUtils

import (
	// "context"
	// "errors"
	// "fmt"
	// "io"
	// "os"

	"minik8s/pkg/k8stype"
	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/api/types/filters"
)

/* imageManager 调用 kubecontainer 提供的PullImage/GetImageRef/
 * ListImages/RemoveImage/ImageStates 方法来保证pod 运行所需要的镜像。
 */
type ImageManager interface {
	PullImage(imageRef string, policy k8stype.ImagePullPolicy) (string, error)
}

type imageManagerImpl struct {
}

var imageManager *imageManagerImpl = nil

func GetImageManager() ImageManager {
	if imageManager == nil {
		imageManager = &imageManagerImpl{}
	}
	return imageManager
}

func (imageMgr *imageManagerImpl) PullImage(imageRef string, policy k8stype.ImagePullPolicy) (string, error) {
	return "", nil
	// ctx := context.Background()

	// if client, err := NewDockerClient(); err != nil {
	// 	return "", err
	// }

	// defer client.Close()

	// switch policy {
	// // 无论如何都要拉取镜像
	// case k8stype.PullAlways:
	// 	// 拉取镜像
	// 	image, err := client.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	// 	// println(imageRef)
	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	file, err := os.Create(os.DevNull)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	defer file.Close()
	// 	io.Copy(file, image)
	// 	defer image.Close()

	// 	imageIDs, err := imageMgr.findLocalImageIDsByImageRef(imageRef)

	// 	if err != nil {
	// 		return "", err
	// 	}

	// 	if len(imageIDs) != 1 {
	// 		return "", errors.New("imageID not found or more than one imageID")
	// 	}

	// 	return imageIDs[0], nil

	// // 如果镜像不存在，那么就拉取镜像
	// case k8stype.PullIfNotPresent, "":
	// 	// 在本地查找镜像
	// 	imageIDs, err := imageMgr.findLocalImageIDsByImageRef(imageRef)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if len(imageIDs) == 0 {
	// 		// 拉取镜像
	// 		image, err := client.ImagePull(ctx, imageRef, types.ImagePullOptions{})

	// 		if err != nil {
	// 			return "", err
	// 		}
	// 		io.Copy(os.Stdout, image)

	// 		defer image.Close()

	// 		imageIDs, err := imageMgr.findLocalImageIDsByImageRef(imageRef)

	// 		if err != nil {
	// 			return "", err
	// 		}

	// 		if len(imageIDs) != 1 {
	// 			fmt.Println(imageIDs)
	// 			return "", errors.New("imageID not found or more than one imageID")
	// 		}

	// 		return imageIDs[0], nil
	// 	} else {
	// 		return imageIDs[0], nil
	// 	}
	// case k8stype.PullNever:
	// 	// 在本地查找镜像
	// 	imageIDs, err := imageMgr.findLocalImageIDsByImageRef(imageRef)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if len(imageIDs) == 0 {
	// 		// 返回一个空的imageID和错误
	// 		return "", errors.New("image not found")
	// 	}
	// 	return imageIDs[0], nil
	// }

	// // 创建一个错误返回
	// err = errors.New("policy not found or not supported")
	// return "", err
}

func (imageMgr *imageManagerImpl) findLocalImage(imageRef string) ([]string, error) {
	// ctx := context.Background()
	// client, err := NewDockerClient()
	// if err != nil {
	// 	// 返回一个空的切片和错误
	// 	return []string{}, err
	// }

	// defer client.Close()
	// // 创建一个过滤器
	// filterArgs := filters.NewArgs()
	// filterArgs.Add("reference", parseImageRef(imageRef))
	// // 查找镜像

	// images, err := client.ImageList(ctx, types.ImageListOptions{
	// 	Filters: filterArgs,
	// })

	// // 创建一个空的切片
	// imageIDs := []string{}

	// // 遍历输出的镜像
	// for _, image := range images {
	// 	imageIDs = append(imageIDs, image.ID)
	// }

	// if err != nil {
	// 	return []string{}, err
	// }

	// return imageIDs, nil
	return []string{""}, nil
}
