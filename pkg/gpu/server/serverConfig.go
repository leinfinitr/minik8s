package server
type TaskServerConfig struct{
	UserName string
	PassWord string
	ServerUri string
	BaseDir string
	RemoteBaseDir string
	TaskName string
	TaskNameSpace string
	OutputPath string
	ErrorPath string
	Partition string
	GPUNum int
	RunCmd []string
}

const (
	submitSbatch = "sbatch -D %s %s"
	statusSbatch = "sacct -j %s --format=JobID,JobName,Partition,State,ExitCode --noheader"
)