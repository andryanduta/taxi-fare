package main

import (
	"log"

	core "github.com/andryanduta/taxi-fare/core"
	"github.com/andryanduta/taxi-fare/fareevaluator"
	fareevaluatorhandler "github.com/andryanduta/taxi-fare/fareevaluator/handler"
	fareevaluatorservice "github.com/andryanduta/taxi-fare/fareevaluator/service"
)

func main() {

	mainConfig, err := core.InitMainConfig()
	if err != nil {
		log.Fatalln(err)
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Init Fare Evaluator
	fareRulesSetting := make(map[fareevaluatorservice.Scope]fareevaluatorservice.FareRule)
	for key, val := range mainConfig.FaseRule {
		scope, ok := fareevaluatorservice.ScopeValue[key]
		if !ok {
			// Skip if scope key rule is unrecognized by Fare Evaluator service
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
