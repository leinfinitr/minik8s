package main

import (
	"encoding/json"
	"flag"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"github.com/mholt/archiver"
	"minik8s/pkg/gpu/server"
	"minik8s/tools/log"
	"net/http"
	"os"
	"strings"
)

var (
	jobDir = "/job"
	remoteJobDir = "job"
	resDir = jobDir + "/res"
)
func main() {
	var jobName string
	var jobNamespace string
	var serverAddr string
	flag.StringVar(&jobName,"jobName","jobName","jobName")
	flag.StringVar(&jobNamespace,"jobNamespace","jobNamespace","jobNamespace")
	flag.StringVar(&serverAddr,"serverAddr",config.APIServerURL(),"serverAddr")
	flag.Parse()
	if jobName == "" || jobNamespace == "" {
		log.ErrorLog("jobName or jobNamespace is empty")
		return
	}
	url := serverAddr + config.JobURI
	url = strings.Replace(url, config.NameSpaceReplace, jobNamespace, -1)
	url = strings.Replace(url, config.NameReplace, jobName, -1)
	res,err := http.Get(url)
	if err != nil {
		log.ErrorLog("GetJob: " + err.Error())
		return
	}
	defer res.Body.Close()
	var job apiObject.Job
	err = json.NewDecoder(res.Body).Decode(&job)
	if err != nil {
		log.ErrorLog("GetJob: " + err.Error())
		return
	}
	log.InfoLog("GetJob: " + job.Metadata.Name)
	serverConfig := &server.TaskServerConfig{
		TaskName: job.Metadata.Name,
		TaskNameSpace: job.Metadata.Namespace,
		ServerUri: serverAddr,
		GPUNum: job.Spec.GPUNum,
		BaseDir: jobDir,
		RemoteBaseDir: remoteJobDir,
		RunCmd: job.Spec.RunCmd,
		UserName: job.Spec.UserName,
		PassWord: job.Spec.PassWord,
		Partition: job.Spec.Partition,
		OutputPath: job.Spec.OutputFile,
		ErrorPath: job.Spec.ErrorFile,
	}
	url = serverAddr + config.JobCodeURI
	url = strings.Replace(url, config.NameSpaceReplace, job.Metadata.Namespace, -1)
	url = strings.Replace(url, config.NameReplace, job.Metadata.Name, -1)
	var jobCode apiObject.JobCode
	res,err = http.Get(url)
	if err != nil {
		log.ErrorLog("GetJobCode: " + err.Error())
		return
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(jobCode)
	if err != nil {
		log.ErrorLog("GetJobCode: " + err.Error())
		return
	}
	jobFileName := "jobzip-"+job.Metadata.UUID
	jobFileNameFull := jobFileName +".zip"
	err = os.RemoveAll(jobDir)
	if err != nil {
		log.ErrorLog("RemoveAll: " + err.Error())
		return
	}
	err = os.MkdirAll(jobDir, 0777)
	if err != nil {
		log.ErrorLog("MkdirAll: " + err.Error())
		return
	}
	err = os.WriteFile(jobDir+"/"+jobFileNameFull,jobCode.UploadContent,0777)
	if err != nil {
		log.ErrorLog("WriteFile: " + err.Error())
		return
	}
	err = archiver.Zip.Open(jobDir+"/"+jobFileNameFull,jobDir)
	if err != nil {
		log.ErrorLog("Open: " + err.Error())
		return
	}
	submitDir := job.Spec.SubmitDir
	idx := strings.LastIndex(submitDir,"/")
	if idx == -1 {
		log.ErrorLog("SubmitDir is not correct")
		return
	}
	submitDir = submitDir[idx+1:]
	if len(submitDir) == 0 {
		log.ErrorLog("SubmitDir is not correct")
		return
	}
	err = os.Rename(serverConfig.BaseDir+"/"+submitDir,serverConfig.BaseDir+"/"+jobFileName)
	if err != nil {
		log.ErrorLog("Rename: " + err.Error())
		return
	}
	serverConfig.BaseDir = serverConfig.BaseDir + "/" + jobFileName
	serverConfig.RemoteBaseDir = serverConfig.RemoteBaseDir + "/" + jobFileName
	ts,err := server.NewTaskServer(serverConfig)
	if err != nil {
		log.ErrorLog("NewTaskServer: " + err.Error())
		return
	}
	ts.RunTaskServer()
}


