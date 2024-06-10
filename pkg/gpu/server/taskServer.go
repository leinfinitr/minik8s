package server

import (
	"fmt"
	"io/fs"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	sshclient "minik8s/pkg/gpu/ssh"
	"minik8s/tools/executor"
	"minik8s/tools/log"
	"minik8s/tools/netRequest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type TaskServer struct {
	config *TaskServerConfig
	client sshclient.SSHClient
}

type JobParsedStatus struct {
	JobID     string `json:"jobID" yaml:"jobID"`
	JobName   string `json:"jobName" yaml:"jobName"`
	Partition string `json:"partition" yaml:"partition"`
	State     string `json:"state" yaml:"state"`
	ExitCode  int    `json:"exitCode" yaml:"exitCode"`
}

func NewTaskServer(config *TaskServerConfig) (*TaskServer, error) {
	log.InfoLog("username: " + config.UserName)
	log.InfoLog("password: " + config.PassWord)
	client, err := sshclient.NewSSHClient(config.UserName, config.PassWord)
	if err != nil {
		return nil, err
	}
	return &TaskServer{
		config: config,
		client: client,
	}, nil
}

func (ts *TaskServer) UploadFiles(filePathList []string) error {
	for _, filePath := range filePathList {
		relativePath, err := filepath.Rel(ts.config.BaseDir, filePath)
		if err != nil {
			return err
		}
		remotePath := path.Join(ts.config.RemoteBaseDir, relativePath)
		err = ts.client.UploadFile(filePath, remotePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ts *TaskServer) InitializeTaskServer() error {
	_, err := ts.client.RmDir(ts.config.RemoteBaseDir)
	if err != nil {
		return err
	}
	_, err = ts.client.MkDir(ts.config.RemoteBaseDir)
	if err != nil {
		return err
	}
	filesNeedUpload := ts.searchFiles()
	err = ts.UploadFiles(filesNeedUpload)
	if err != nil {
		log.ErrorLog("UploadFiles: " + err.Error())
		return err
	}
	scriptContent := ts.GenerateScript()
	_,err = ts.client.WriteFile(path.Join(ts.config.RemoteBaseDir, ts.TaskPath()), scriptContent)
	if err != nil {
		log.ErrorLog("WriteFile: " + err.Error())
		return err
	}
	return nil
}

func (ts *TaskServer) GenerateScript() string {
	scriptContent := "#!/bin/bash\n"
	scriptContent += fmt.Sprintf("#SBATCH --job-name=%s\n", ts.config.TaskName)
	scriptContent += fmt.Sprintf("#SBATCH --output=%s\n", ts.config.OutputPath)
	scriptContent += fmt.Sprintf("#SBATCH --error=%s\n", ts.config.ErrorPath)
	scriptContent += fmt.Sprintf("#SBATCH --partition=%s\n", ts.config.Partition)
	scriptContent += fmt.Sprintf("#SBATCH --gres=gpu:%d\n", ts.config.GPUNum)
	scriptContent += "#SBATCH -N 1\n"
	scriptContent += ts.FormatRunCmd()
	
	return scriptContent
}

func (ts *TaskServer) FormatRunCmd() string {
	totalContent := ""
	for _, v := range ts.config.RunCmd {
		totalContent += v + "\n"
	}
	return totalContent
}
func (ts *TaskServer) searchFiles() []string {
	filePathList := make([]string, 0)
	filepath.WalkDir(ts.config.BaseDir, func(path string, entry fs.DirEntry, err error) error {
		if !entry.IsDir() {
			filename := entry.Name()
			if filepath.Ext(filename) == ".cu" {
				filePathList = append(filePathList, path)
			}
		}
		return nil
	})
	return filePathList
}

func (ts *TaskServer) TaskPath() string {
	return ts.config.TaskName + ".slurm"
}

func (ts *TaskServer) CheckStatus(taskId string) (*JobParsedStatus, error) {
	cmd := fmt.Sprintf(statusSbatch, taskId)
	res, err := ts.client.RunCmd(cmd)
	if err != nil {
		log.ErrorLog("RunCmd: " + err.Error())
		return nil, err
	}
	var status JobParsedStatus
	num, err := fmt.Sscanf(res, "%s %s %s %d", &status.JobID, &status.Partition, &status.State, &status.ExitCode)
	if err != nil || num != 4 {
		log.ErrorLog("RunCmd: " + err.Error())
		return nil, err
	}
	return &status, nil
}

func (ts *TaskServer) CheckAndUpdateStatus(taskId string) int {
	status, err := ts.CheckStatus(taskId)
	if err != nil {
		log.ErrorLog("CheckStatus: " + err.Error())
		return 2
	}
	if status != nil {
		jStatus := &apiObject.JobStatus{
			JobID:     taskId,
			Partition: status.Partition,
			State:     status.State,
			ExitCode:  status.ExitCode,
		}
		url := ts.config.ServerUri + config.JobStatusURI
		url = strings.Replace(url, config.NameSpaceReplace, ts.config.TaskNameSpace, -1)
		url = strings.Replace(url, config.NameReplace, ts.config.TaskName, -1)
		code, _, err := netRequest.PutRequestByTarget(url, jStatus)
		if err != nil {
			log.ErrorLog("PostRequestByTarget: " + err.Error())
			return 2
		}
		if code != 200 {
			log.ErrorLog("PostRequestByTarget: code is not 200")
			return 2
		}
		return 1
	}
	return 0
}

func (ts *TaskServer)ProcessResult() error{
	remoteOutputPath := path.Join(ts.config.RemoteBaseDir, ts.config.OutputPath)
	remoteErrorPath := path.Join(ts.config.RemoteBaseDir, ts.config.ErrorPath)
	log.InfoLog("ProcessResult: " + remoteOutputPath)
	log.InfoLog("ProcessResult: " + remoteErrorPath)
	localOutputPath := path.Join(ts.config.BaseDir, ts.config.OutputPath)
	err := ts.client.DownloadFile(remoteOutputPath, localOutputPath)
	if err != nil {
		log.ErrorLog("DownloadFile: " + err.Error())
		return err
	}
	localErrorPath := path.Join(ts.config.BaseDir, ts.config.ErrorPath)
	err = ts.client.DownloadFile(remoteErrorPath, localErrorPath)
	if err != nil {
		log.ErrorLog("DownloadFile: " + err.Error())
		return err
	}
	outputContent,err := os.ReadFile(localOutputPath)
	if err != nil {
		log.ErrorLog("ReadFile: " + err.Error())
		return err
	}
	errorContent,err := os.ReadFile(localErrorPath)
	if err != nil {
		log.ErrorLog("ReadFile: " + err.Error())
		return err
	}
	var jobCode apiObject.JobCode
	jobCode.OutputContent = outputContent
	jobCode.ErrorContent = errorContent
	log.InfoLog("ProcessResult: " + string(outputContent))
	log.InfoLog("ProcessResult: " + string(errorContent))
	url := ts.config.ServerUri + config.JobCodeURI
	url = strings.Replace(url, config.NameSpaceReplace, ts.config.TaskNameSpace, -1)
	url = strings.Replace(url, config.NameReplace, ts.config.TaskName, -1)
	code, _, err := netRequest.PutRequestByTarget(url, jobCode)
	if err != nil {
		log.ErrorLog("PostRequestByTarget: " + err.Error())
		return err
	}
	if code != 200 {
		log.ErrorLog("PostRequestByTarget: code is not 200")
		return err
	}
	_,err = ts.client.RmDir(ts.config.RemoteBaseDir)
	if err != nil {
		log.ErrorLog("RmDir: " + err.Error())
		return err
	}
	return nil
}

func (ts *TaskServer) RunTaskServer() {
	err := ts.InitializeTaskServer()
	if err != nil {
		log.ErrorLog("InitializeTaskServer: " + err.Error())
		return
	}
	//提交作业
	cmd := fmt.Sprintf(submitSbatch, ts.config.RemoteBaseDir, path.Join(ts.config.RemoteBaseDir,ts.TaskPath()))
	log.InfoLog("RunTaskServer: " + cmd)
	res, err := ts.client.RunCmd(cmd)
	if err != nil {
		log.ErrorLog("RunCmd: " + err.Error())
		return
	}
	var taskId string
	num, err := fmt.Sscanf(res, "Submitted batch job %s", &taskId)
	if err != nil || num != 1 {
		log.ErrorLog("RunCmd: " + err.Error())
		return
	}
	log.InfoLog("RunTaskServer: " + taskId)
	//监控作业
	checkInPeriod := func() bool {
		switch ts.CheckAndUpdateStatus(taskId) {
		case 0:
			log.InfoLog("Task processing")
			return false
		case 1:
			time.Sleep(2 * time.Second)
			log.InfoLog("Task finished")
			err := ts.ProcessResult()
			if err != nil {
				log.ErrorLog("ProcessResult: " + err.Error())
				return false
			}
			return true
		case 2:
			log.ErrorLog("Task failed")
			return false
		}
		return false
	}
	executor.CondExecuteInPeriod(0*time.Second, []time.Duration{5 * time.Second}, checkInPeriod)
}
