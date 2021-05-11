package main

import (
	"flag"
	"fmt"
	"github.com/proxy-go/pkg/config"
	"github.com/proxy-go/pkg/log"
	"github.com/proxy-go/pkg/proxy"
	"net/http"
	"os"
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

	if err = http.ListenAndServe(fmt.Sprintf(":%d", cmdInput.Port), mux); err != nil {
		panic(err)
	}
}
