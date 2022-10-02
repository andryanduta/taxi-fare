package handler

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/andryanduta/taxi-fare/fareevaluator"
)

const (
	timeFormat   = "15:04:05.000"
	timeInterval = 300
	validPattern = `^\d\d:\d\d:\d\d.\d\d\d\s\d+\.?\d*$`
	maxMileage   = float64(99999999.9)
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
	type IndexMileage struct {
		Index       int
		MileageDiff float64
	}
	var idxMileages []IndexMileage
	var mileageBefore float64
	for i, input := range inputs {
		distanceMeter, err := constructInput(input)
		if err != nil {
			log.Fatalf("[Handler Fare Evaluator] error constructInput, input: %s, err: %+v\n", input, err)
		}
		distanceMeters = append(distanceMeters, distanceMeter)
		curr := distanceMeter.Mileage - mileageBefore
		idxMileages = append(idxMileages, IndexMileage{
			Index:       i,
			MileageDiff: curr,
		})
		mileageBefore = distanceMeter.Mileage
	}

	res, err := h.adEvaluatorSvc.CalculateFare(distanceMeters)
	if err != nil {
		log.Fatalf("[Handler Fare Evaluator] error CalculateFare, distanceMeters: %+v, err: %+v\n", distanceMeters, err)
	}

	sort.Slice(idxMileages, func(i, j int) bool {
		return idxMileages[i].MileageDiff > idxMileages[j].MileageDiff
	})

	// display all of the input data with mileage difference compared to previous
	// data, and order it from highest to lowest
	log.Println(res)
	for _, data := range idxMileages {
		log.Println(distanceMeters[data.Index].ElapsedTime, distanceMeters[data.Index].Mileage, data.MileageDiff)
	}
}

func isMileageValid(current float64, before float64) bool {
	if current == 0 {
		return false
	}
	if current >= maxMileage {
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
