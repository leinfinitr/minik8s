package sshclient

import "os"

var (
	SSHUserName = ""
	SSHPassword = ""
	SSHPort     = 22
)

const (
	sshLoginAddr = "pilogin.hpc.sjtu.edu.cn"
	sshDataAddr = "data.hpc.sjtu.edu.cn"
)

const (
	ChDirCmd          = "cd "
	RmDirCmd          = "rm -rf "
	MkDirCmd          = "mkdir -p "
	RdFileCmd         = "cat "
	WriteFileCmd      = "echo "
	AppendFileCommand = "echo "
	RmFileCommand     = "rm -rf "
)

func init() {
	// 从环境变量中读取用户名和密码
	SSHUserName = os.Getenv("GPU_SSH_USERNAME")
	SSHPassword = os.Getenv("GPU_SSH_PASSWORD")
}
