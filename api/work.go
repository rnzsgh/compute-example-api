package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	log "github.com/golang/glog"
	"github.com/google/uuid"
	"github.com/rnzsgh/compute-example-api/cloud"
	"github.com/rnzsgh/compute-example-api/model"
)

func init() {
	http.HandleFunc("/submit", WorkSubmit)
}

func WorkSubmit(w http.ResponseWriter, r *http.Request) {
	response := &response{}
	w.Header().Set("Content-Type", "text/plain")

	msg := &model.JobMessage{
		ObjectKey: "someObjectKeyInS3",
		Type:      "extractInsight",
		RequestId: uuid.New().String(),
	}

	msgRaw, _ := json.Marshal(msg)

	if messageId, err := cloud.SqsSendMessage(r.Context(), os.Getenv("SQS_JOB_QUEUE_URL"), string(msgRaw)); err != nil {
		response.Message = "Error"
		w.WriteHeader(http.StatusInternalServerError)
		log.Error(err)
	} else {
		log.Info("Submitted message to queue: %s - message id: %s", os.Getenv("SQS_JOB_QUEUE_URL"), messageId)
		response.Message = "Accepted"
		w.WriteHeader(http.StatusAccepted)
	}

	out, _ := json.Marshal(response)
	io.WriteString(w, string(out))
}
