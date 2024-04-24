package image

/* imageManager 调用 kubecontainer 提供的PullImage/GetImageRef/
 * ListImages/RemoveImage/ImageStates 方法来保证pod 运行所需要的镜像。
 */
type ImageManager struct {
}
