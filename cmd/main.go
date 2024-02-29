package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"

	"file-diff-finder/cache"
	"file-diff-finder/config"
	df "file-diff-finder/diff_finder"
	"file-diff-finder/router"
)

func main() {
	cnf := config.InitConfig()
	cacheInstance := cache.InitCache()
	originalFile := df.NewFileInfo(cnf.FileContent, cnf.FileVersion)
	logger := logrus.New()
	rhs := df.NewFileDiffFinder(originalFile)

	rt := router.NewRouter(logger, cacheInstance, rhs)
	logger.Info("Service is initializing...")

	var err error
	if !errors.Is(err, http.ListenAndServe(":"+strconv.Itoa(cnf.Port), rt)) {
		log.Fatalf("error initializing router: %v", err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
}
