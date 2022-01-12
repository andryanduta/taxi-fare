package handler

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/andryanduta/taxi-fare/fareevaluator"
)

const (
	timeFormat   = "15:04:05.000"
	timeInterval = 300
	validPattern = `^\d\d:\d\d:\d\d.\d\d\d\s\d+\.?\d*$`
)

type handler struct {
	adEvaluatorSvc fareevaluator.Service
}

func New(svc fareevaluator.Service) *handler {
	handler := &handler{
		adEvaluatorSvc: svc,
	}

	return handler
}

func (h *handler) HandleFareEvaluator() {
	reader := bufio.NewReader(os.Stdin)
	inputs := make([]string, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			log.Println("[Handler Fare Evaluator] error ReadLine", err)
			break
		}

		strInput := string(line)
		if strInput == "" {
			break
		}

		inputs = append(inputs, strInput)
	}
	err := validateInputs(inputs)
	if err != nil {
		log.Fatalf("[Handler Fare Evaluator] bad request input, inputs: %s, err: %+v\n", inputs, err)
	}
	var distanceMeters []fareevaluator.DistanceMeter
	for _, input := range inputs {
		distanceMeter, err := constructInput(input)
		if err != nil {
			log.Fatalf("[Handler Fare Evaluator] error constructInput, input: %s, err: %+v\n", input, err)
		}
		distanceMeters = append(distanceMeters, distanceMeter)
	}

	res, err := h.adEvaluatorSvc.CalculateFare(distanceMeters)
	if err != nil {
		log.Fatalf("[Handler Fare Evaluator] error CalculateFare, distanceMeters: %+v, err: %+v\n", distanceMeters, err)
	}

	log.Println(res)
}

func isMileageValid(current float64, before float64) bool {
	if current == 0 {
		return false
	}
	return before <= current
}

func isElapsedTimeValid(current time.Time, before time.Time) bool {
	if current.Before(before) {
		return false
	}

	return current.Sub(before) <= time.Duration(timeInterval)*time.Second
}

func validateInputs(inputs []string) error {
	// validate lines data
	if !(len(inputs) >= 2) {
		return ErrInvalidDataCount
	}

	var timeBefore time.Time
	var mileageBefore float64

	for i, inp := range inputs {

		// validate input format pattern
		rgx := regexp.MustCompile(validPattern)
		if !rgx.Match([]byte(inp)) {
			return ErrInvalidFormat
		}

		splitString := strings.Split(inp, " ")
		elapsedTime, err := time.Parse(timeFormat, splitString[0])
		if err != nil {
			return err
		}

		mileage, err := strconv.ParseFloat(splitString[1], 64)
		if err != nil {
			return err
		}

		if i > 0 {
			// validate mileage & input time before vs current
			if !isMileageValid(mileage, mileageBefore) {
				return ErrInvalidMileage
			}
			if !isElapsedTimeValid(elapsedTime, timeBefore) {
				return ErrInvalidElapsedTime
			}
		}

		timeBefore = elapsedTime
		mileageBefore = mileage
	}
	return nil
}

func constructInput(line string) (distanceMeter fareevaluator.DistanceMeter, err error) {
	splitLine := strings.Split(line, " ")

	mileage, err := strconv.ParseFloat(splitLine[1], 64)
	if err != nil {
		return distanceMeter, err
	}

	distanceMeter.ElapsedTime = splitLine[0]
	distanceMeter.Mileage = mileage
	return distanceMeter, nil
}
