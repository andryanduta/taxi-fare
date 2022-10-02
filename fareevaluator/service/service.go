package service

import (
	"math"

	"github.com/andryanduta/taxi-fare/fareevaluator"
)

type Service struct {
	config Config
}

type Config struct {
	fareRule map[Scope]FareRule
}

// Scope denotes a shared settings.
type Scope int

// The following are the known Scopes.
const (
	_ Scope = iota
	ScopeFareBase
	ScopeFareUpTo
	ScopeFareOver
)

// String returns string representation of Scope.
func (sc Scope) String() string { return ScopeName[sc] }

// ScopeName defines all known Scopes and their string representation, while ScopeValue
// is the reverse-mapping of ScopeName. Registering new Scope is done by adding new entry in
// both variables.
var (
	ScopeName = map[Scope]string{
		ScopeFareBase: "Base",
		ScopeFareUpTo: "UpTo",
		ScopeFareOver: "Over",
	}

	ScopeValue = map[string]Scope{
		ScopeName[ScopeFareBase]: ScopeFareBase,
		ScopeName[ScopeFareUpTo]: ScopeFareUpTo,
		ScopeName[ScopeFareOver]: ScopeFareOver,
	}
)

type FareRule struct {
	Price             float64
	Distance          float64
	DistanceThreshold float64
}

type Option func(*Service)

func WithConfig(fareRule map[Scope]FareRule) Option {
	return Option(func(s *Service) {
		s.config = Config{
			fareRule: fareRule,
		}
	})
}

func New(options ...Option) *Service {
	service := &Service{}

	for _, opt := range options {
		opt(service)
	}

	return service
}

func (s *Service) CalculateFare(distanceMeters []fareevaluator.DistanceMeter) (taxiFareRounded int64, err error) {
	var taxiFare float64
	if len(distanceMeters) < 2 {
		return 0, fareevaluator.ErrDistanceMeterArraySize
	}

	taxiFare = s.config.fareRule[ScopeFareBase].Price

	// base calculation
	mileage := (distanceMeters[len(distanceMeters)-1].Mileage - distanceMeters[0].Mileage)
	if mileage <= s.config.fareRule[ScopeFareBase].Distance {
		return int64(taxiFare), nil
	}

	// upto calculation
	mileage = mileage - s.config.fareRule[ScopeFareBase].Distance
	remainderDistanceFareRule := s.config.fareRule[ScopeFareUpTo].DistanceThreshold - s.config.fareRule[ScopeFareBase].Distance
	if mileage <= remainderDistanceFareRule {
		taxiFare = taxiFare + (s.config.fareRule[ScopeFareUpTo].Price * (mileage / s.config.fareRule[ScopeFareUpTo].Distance))
		taxiFareRounded = int64(math.Round((taxiFare)))
		return taxiFareRounded, nil
	}
	taxiFare = taxiFare + (s.config.fareRule[ScopeFareUpTo].Price * (remainderDistanceFareRule / s.config.fareRule[ScopeFareUpTo].Distance))

	// over calculation
	taxiFare = taxiFare + (s.config.fareRule[ScopeFareOver].Price * ((mileage - remainderDistanceFareRule) / s.config.fareRule[ScopeFareOver].Distance))
	taxiFareRounded = int64(math.Round((taxiFare)))
	return taxiFareRounded, nil
}
