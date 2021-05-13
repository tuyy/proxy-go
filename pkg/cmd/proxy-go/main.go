package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/tuyy/proxy-go/pkg/config"
	"github.com/tuyy/proxy-go/pkg/log"
	"github.com/tuyy/proxy-go/pkg/proxy"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type CmdInput struct {
	Port         int
	ConfFilePath string
}

var cmdInput CmdInput

func init() {
	flag.IntVar(&cmdInput.Port, "p", -1, "input server port")
	flag.StringVar(&cmdInput.ConfFilePath, "c", "", "input conf file path")
	flag.Parse()

	jsonConfData, err := os.ReadFile(cmdInput.ConfFilePath)
	if err != nil {
		panic(err)
	}

	config.Init(jsonConfData)
	log.Init(config.Loggers())
}

func main() {
	mux, err := proxy.MakeServeMux(config.DefaultDestUrl(), config.URLMappings())
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cmdInput.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	closeCh := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Info("HTTP server Shutdown: %v", err)
		}
		close(closeCh)
	}()

	log.Info("START PROXY SERVER PORT:%d CONF:%s", cmdInput.Port, cmdInput.ConfFilePath)

	if err = srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error("HTTP server ListenAndServe: %v", err)
	}

	log.Info("STOP PROXY SERVER PORT:%d CONF:%s", cmdInput.Port, cmdInput.ConfFilePath)

	<-closeCh
}
