package sshclient

import (
	"minik8s/tools/log"
	"os"
	"path/filepath"

	"github.com/melbahja/goph"
)

type SSHClient interface {
	RunCmd(cmd string) (string, error)
	RunCmds(cmds []string) (string, error)

	ChDir(path string) (string, error)
	RmDir(path string) (string, error)
	MkDir(path string) (string, error)

	RdFile(path string) (string, error)
	WriteFile(path string, content string) (string, error)
	AppendFile(path string, content string) (string, error)
	RmFile(path string) (string, error)

	UploadFile(localPath string, remotePath string) error
	DownloadFile(remotePath string, localPath string) error

	UploadDir(localPath string, remotePath string) error
}

type sshClient struct {
	client   *goph.Client
	userName string
	passWord string
}

func NewSSHClient(userName, pwd string) (SSHClient, error) {
	cli, err := goph.NewUnknown(userName, sshLoginAddr, goph.Password(pwd))
	if err != nil {
		return nil, err
	}

	return &sshClient{
		client:   cli,
		userName: userName,
		passWord: pwd,
	}, nil
}

func (sshCli *sshClient) MkDir(path string) (string, error) {
	output, err := sshCli.client.Run(MkDirCmd + path)
	return string(output), err
}

func (sshCli *sshClient) RmDir(path string) (string, error) {
	output, err := sshCli.client.Run(RmDirCmd + path)
	return string(output), err
}

func (sshCli *sshClient) ChDir(path string) (string, error) {
	output, err := sshCli.client.Run(ChDirCmd + path)
	return string(output), err
}

func (sshCli *sshClient) RdFile(path string) (string, error) {
	output, err := sshCli.client.Run(RdFileCmd + path)
	return string(output), err
}

func (sshCli *sshClient) WriteFile(path string, content string) (string, error) {
	wCmd := WriteFileCmd + "\"" + content + "\" > " + path
	output, err := sshCli.client.Run(wCmd)
	return string(output), err
}

func (sshCli *sshClient) AppendFile(path string, content string) (string, error) {
	appCmd := AppendFileCommand + "\"" + content + "\" >> " + path
	output, err := sshCli.client.Run(appCmd)
	return string(output), err
}

func (sshCli *sshClient) RmFile(path string) (string, error) {
	output, err := sshCli.client.Run(RmFileCommand + path)
	return string(output), err
}

func (sshCli *sshClient) RunCmd(cmd string) (string, error) {
	output, err := sshCli.client.Run(cmd)
	return string(output), err
}

func (sshCli *sshClient) RunCmds(cmds []string) (string, error) {
	var output string
	var err error
	for _, cmd := range cmds {
		curOutput, err := sshCli.client.Run(cmd)
		if err != nil {
			return output, err
		}
		output += string(curOutput)
	}
	return output, err
}

func (sshCli *sshClient) UploadFile(localPath string, remotePath string) error {
	err := sshCli.client.Upload(localPath, remotePath)
	return err
}

func (sshCli *sshClient) DownloadFile(remotePath string, localPath string) error {
	log.InfoLog("remoteFIlePath:"+ remotePath)
	log.InfoLog("localFilePath:"+ localPath)
	err := sshCli.client.Download(remotePath, localPath)
	return err
}


func (sshCli *sshClient) UploadDir(localPath string, remotePath string) error {
	err := filepath.Walk(localPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		relativePath, err := filepath.Rel(localPath, path)
		if err != nil {
			return err
		}
		remoteAbsPath := filepath.Join(remotePath, relativePath)
		err = sshCli.client.Upload(path, remoteAbsPath)

		if err != nil {
			log.ErrorLog("Upload file failed:"+ remoteAbsPath)
		} else {
			log.InfoLog("Upload file success:"+ remoteAbsPath)
		}

		return nil
	})

	return err
}
