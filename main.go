package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	"github.com/rnzsgh/compute-example-api/api"
	"github.com/rnzsgh/compute-example-api/listen"
)

func init() {
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
}

func main() {

	sigs := make(chan os.Signal, 1)
	running := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	srv := api.RunHttpServer()

	stopListening := listen.ListenForCompletedJobs()

	go func() {
		<-sigs
		running <- false
	}()

	<-running

	stopListening()

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Errorf("Unable to shutdown http server: %v", err)
	}

	log.Flush()
}
