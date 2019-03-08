package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/golang/glog"
	"github.com/rnzsgh/compute-example-api/api"
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

	go func() {
		<-sigs
		running <- false
	}()

	<-running

	if err := srv.Shutdown(context.TODO()); err != nil {
		log.Errorf("Unable to shutdown http server: %v", err)
	}

	log.Flush()
}
