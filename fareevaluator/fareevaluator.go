package fareevaluator

type DistanceMeter struct {
	Mileage     float64
	ElapsedTime string
}

// Service is the interface for fare evaluator related
//
//go:generate mockgen -destination mockservice/mock_service.go -package mockservice github.com/andryanduta/taxi-fare/fareevaluator Service

type Service interface {
	CalculateFare(distanceMeters []DistanceMeter) (float64, error)
}

var defaultService Service

// Init sets the default Fare Evaluator Service.
func Init(svc Service) {
	defaultService = svc
}
