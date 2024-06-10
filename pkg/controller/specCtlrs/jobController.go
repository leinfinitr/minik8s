package specctlrs

import (
	"encoding/json"
	"minik8s/pkg/apiObject"
	"minik8s/pkg/config"
	"minik8s/tools/executor"
	"minik8s/tools/log"
	netRequest "minik8s/tools/netRequest"
	"net/http"
	"strings"
	"time"
)

type JobController interface {
	Run()
}
type JobControllerImpl struct {
}

var (
	JobControllerDelay   = 3 * time.Second
	JobControllerTimeGap = []time.Duration{14 * time.Second}
	ServerImage          = "jackhel0/task-server:latest"
)

func NewJobController() (JobController, error) {
	return &JobControllerImpl{}, nil
}
func (jc *JobControllerImpl) syncJob() {
	url := config.APIServerURL() + config.GlobalJobURI
	res, err := http.Get(url)
	if err != nil {
		log.ErrorLog("syncJob: " + err.Error())
		return
	}
	defer res.Body.Close()
	var jobs []apiObject.Job
	err = json.NewDecoder(res.Body).Decode(&jobs)
	if err != nil {
		log.ErrorLog("syncJob: " + err.Error())
		return
	}
	var emptyStatus = apiObject.JobStatus{}
	for _, job := range jobs {
		log.InfoLog("syncJob: " + job.Metadata.Name)
		//未处理的Job没有Status信息
		if job.Status == emptyStatus {
			jc.CreateJob(job)
		}
	}

}
func (jc *JobControllerImpl) CreateJob(job apiObject.Job) {
	// cmd := []string{"/bin/server", "-jobNname", job.Metadata.Name, "-jobNamespace", job.Metadata.Namespace,"-serverAddr",config.APIServerURL()}
	cmd := []string{"/bin/sh", "-c", `echo "nameserver 223.5.5.5" > /etc/resolv.conf &&
    /bin/server -jobName ` + job.Metadata.Name + ` -jobNamespace ` + job.Metadata.Namespace + ` -serverAddr ` + config.APIServerURL()}
	pod := &apiObject.Pod{
		TypeMeta: apiObject.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		Metadata: apiObject.ObjectMeta{
			Name:      job.Metadata.Name,
			Namespace: job.Metadata.Namespace,
		},
		Spec: apiObject.PodSpec{
			Containers: []apiObject.Container{
				{
					Name:    "server" + job.Metadata.UUID,
					Image:   ServerImage,
					Command: cmd,
				},
			},
			RestartPolicy: "Always",
		},
	}
	url := config.APIServerURL() + config.PodsURI
	url = strings.Replace(url, config.NameSpaceReplace, job.Metadata.Namespace, -1)
	code, _, err := netRequest.PostRequestByTarget(url, pod)
	if err != nil {
		log.ErrorLog("CreateJob: " + err.Error())
		return
	}
	if code != http.StatusCreated {
		log.ErrorLog("CreateJob: " + "Create Job Failed")
		return
	}
	log.InfoLog("CreateJob: " + "Create Job Success")
}
func (jc *JobControllerImpl) Run() {
	// 定期执行
	executor.ExecuteInPeriod(JobControllerDelay, JobControllerTimeGap, jc.syncJob)
}
