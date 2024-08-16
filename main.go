package main

import (
	"io"
	"os"

	core "github.com/andryanduta/taxi-fare/core"
	"github.com/andryanduta/taxi-fare/fareevaluator"
	fareevaluatorhandler "github.com/andryanduta/taxi-fare/fareevaluator/handler"
	fareevaluatorservice "github.com/andryanduta/taxi-fare/fareevaluator/service"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func setupLogger() {
	// Open a file for logging
	logFile, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	multi := io.MultiWriter(os.Stdout, os.Stderr, logFile)

	log.Logger = zerolog.New(multi).With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)
}

func init() {
	setupLogger()
}

func main() {

	mainConfig, err := core.InitMainConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("[main] error InitMainConfig")
	}

	// Init Fare Evaluator
	fareRulesSetting := make(map[fareevaluatorservice.Scope]fareevaluatorservice.FareRule)
	for key, val := range mainConfig.FaseRule {
		scope, ok := fareevaluatorservice.ScopeValue[key]
		if !ok {
			// Skip if scope key rule is unrecognized by Fare Evaluator service
			log.Debug().Str("key", key).Msg("[main] scope key is unrecognized")
			continue
		}

		fareRulesSetting[scope] = fareevaluatorservice.FareRule{
			Price:             val.Price,
			Distance:          val.Distance,
			DistanceThreshold: val.DistanceThreshold,
		}
	}

	var svcOptions []fareevaluatorservice.Option
	svcOptions = append(svcOptions, fareevaluatorservice.WithConfig(fareRulesSetting))
	fareEvaluatorSvc := fareevaluatorservice.New(
		svcOptions...,
	)
	fareevaluator.Init(fareEvaluatorSvc)

	fareHandler := fareevaluatorhandler.New(fareEvaluatorSvc)

	fareHandler.HandleFareEvaluator()
}
