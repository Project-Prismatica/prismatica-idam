package main

import (
	"flag"

	"github.com/jamiealquiza/envy"
	log "github.com/sirupsen/logrus"

	"github.com/Project-Prismatica/prismatica-idam"
)

type programArguments struct {
	Bind string
	LogDebug, LogVerbose bool
}

const (
	idamEnvVarPrefix = "IDAM"
)

var (
	runtimeArguments = populateProgramArguments()
)

func populateProgramArguments() (args programArguments) {

	flag.StringVar(&args.Bind,
		"bind",
		"127.0.0.1:8080",
		"port for core to bind the auth service",
	)

	flag.BoolVar(&args.LogDebug,"debug", false,
		"use debug log level")

	flag.BoolVar(&args.LogVerbose,"verbose", false,
		"use verbose log level")


	envy.Parse(idamEnvVarPrefix)
	flag.Parse()

	return
}

func configureLogging() {

	if runtimeArguments.LogDebug {
		log.SetLevel(log.DebugLevel)

	} else if runtimeArguments.LogVerbose {
		log.SetLevel(log.InfoLevel)

	} else {
		log.SetLevel(log.WarnLevel)
	}

}

func main() {
	configureLogging()

	log.Info("idam server starting")

	prismatica_idam.RunAmbassadorAuthenticationService(runtimeArguments.Bind)

	log.Info("idam server shutting down")
}
