package handler

import (
	"os"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"

	"minik8s/pkg/apiObject"
	"minik8s/tools/log"
)

// BindEvent 绑定事件
func BindEvent(c *gin.Context) {
	log.DebugLog("BindEvent")
	// 从请求中获取事件内容
	var event apiObject.ServerlessEvent
	err := c.ShouldBindJSON(&event)
	if err != nil {
		log.ErrorLog("BindEvent: " + err.Error())
		c.JSON(400, err.Error())
		return
	}
	// 根据事件类型执行相应操作
	switch event.Type {
	case apiObject.ServerlessEventTypeTime:
		bindTimeEvent(event)
	case apiObject.ServerlessEventTypeFile:
		bindFileEvent(event)
	default:
		log.ErrorLog("BindEvent: Unknown event type.")
		c.JSON(400, "Unknown event type.")
	}
}

// bindTimeEvent 绑定时间事件
//
//	在给定时间长度后触发Serverless Function
func bindTimeEvent(event apiObject.ServerlessEvent) {
	log.DebugLog("BindTimeEvent")
	// 解析事件参数
	duration, err := time.ParseDuration(event.Params)
	if err != nil {
		log.ErrorLog("BindTimeEvent: " + err.Error())
		return
	}
	// 开启一个线程，等待一段时间后触发Serverless Function
	go func() {
		time.Sleep(duration)
		RunFunction(event.Name, event.Args)
	}()
}

// bindFileEvent 绑定文件事件
//
//	当文件内容发生变化时触发Serverless Function
func bindFileEvent(serverlessEvent apiObject.ServerlessEvent) {
	log.DebugLog("BindFileEvent")
	// 监视文件
	fileName := serverlessEvent.Params
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.ErrorLog("failed to create a new watcher: " + err.Error())
		os.Exit(1)
	}
	defer func(watcher *fsnotify.Watcher) {
		err := watcher.Close()
		if err != nil {
			log.ErrorLog("failed to close the watcher: " + err.Error())
			os.Exit(1)
		}
	}(watcher)

	// 确保路径是绝对路径，以避免监视路径解析上的问题
	absPath, err := filepath.Abs(fileName)
	if err != nil {
		log.ErrorLog("failed to get the absolute path of the file: " + err.Error())
		os.Exit(1)
	}

	// 添加文件到监视列表
	if err := watcher.Add(absPath); err != nil {
		log.ErrorLog("failed to add the file to the watcher: " + err.Error())
		os.Exit(1)
	}
	log.InfoLog("Now watching file: " + absPath)

	// 开启一个线程，监视文件变化
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// 当文件被修改时调用Serverless Function
				if event.Op&fsnotify.Write == fsnotify.Write {
					RunFunction(serverlessEvent.Name, serverlessEvent.Args)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.ErrorLog("error:" + err.Error())
			}
		}
	}()

	// 阻塞，直到done通道关闭（通常不会关闭，除非主动结束监控）
	<-done
}
