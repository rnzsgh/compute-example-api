package listen

import (
	"encoding/json"
	"os"
	"time"

	log "github.com/golang/glog"
	"github.com/rnzsgh/compute-example-api/cloud"
	"github.com/rnzsgh/compute-example-api/model"
)

func ListenForCompletedJobs() func() {

	running := true

	go func() {
		for running {
			msg, receipt, err := cloud.SqsReceiveOneMessage(os.Getenv("SQS_JOB_COMPLETED_QUEUE_URL"), 30, 5)
			if err != nil {
				log.Error(err)
				time.Sleep(1 * time.Second)
				continue
			}

			if len(msg) == 0 {
				continue
			}

			job := &model.JobMessage{}
			if err = json.Unmarshal([]byte(msg), job); err != nil {
				log.Errorf("Unable to unmarshal from queue: %s - reason %v", os.Getenv("SQS_JOB_COMPLETED_QUEUE_URL"), err)
			} else {
				log.Infof("Job completed - key: %s - type: %s - request: %s", job.ObjectKey, job.Type, job.RequestId)
			}

			if err = cloud.SqsDeleteMessage(os.Getenv("SQS_JOB_COMPLETED_QUEUE_URL"), receipt); err != nil {
				log.Error(err)
			}
		}
	}()

	return func() {
		running = false
	}
}
