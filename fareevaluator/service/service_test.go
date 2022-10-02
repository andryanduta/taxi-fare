package service

import (
	"reflect"
	"testing"

	"github.com/andryanduta/taxi-fare/fareevaluator"
)

func TestService_CalculateFare(t *testing.T) {
	type fields struct {
		config Config
	}
	type args struct {
		distanceMeters []fareevaluator.DistanceMeter
	}

	defaultFareRule := map[Scope]FareRule{
		ScopeFareBase: FareRule{
			Price:    400,
			Distance: 1000,
		},
		ScopeFareUpTo: FareRule{
			Price:             40,
			Distance:          400,
			DistanceThreshold: 10000,
		},
		ScopeFareOver: FareRule{
			Price:             40,
			Distance:          350,
			DistanceThreshold: 10000,
		},
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantTaxiFare int64
		wantErr      bool
	}{
		{
			name: "Case success. Fare succesfully calculated",
			fields: fields{
				config: Config{
					fareRule: defaultFareRule,
				},
			},
			args: args{
				distanceMeters: []fareevaluator.DistanceMeter{
					{
						Mileage:     0,
						ElapsedTime: "00:00:00.000",
					},
					{
						Mileage:     700,
						ElapsedTime: "00:03:00.123",
					},
				},
			},
			wantTaxiFare: 400,
			wantErr:      false,
		},
		{
			name: "Case success - UpTo. Fare succesfully calculated",
			fields: fields{
				config: Config{
					fareRule: defaultFareRule,
				},
			},
			args: args{
				distanceMeters: []fareevaluator.DistanceMeter{
					{
						Mileage:     0,
						ElapsedTime: "00:00:00.000",
					},
					{
						Mileage:     1100,
						ElapsedTime: "00:03:00.123",
					},
				},
			},
			wantTaxiFare: 410,
			wantErr:      false,
		},
		{
			name: "Case success - Over. Fare succesfully calculated",
			fields: fields{
				config: Config{
					fareRule: defaultFareRule,
				},
			},
			args: args{
				distanceMeters: []fareevaluator.DistanceMeter{
					{
						Mileage:     0,
						ElapsedTime: "00:00:00.000",
					},
					{
						Mileage:     32000,
						ElapsedTime: "00:03:00.123",
					},
				},
			},
			wantTaxiFare: 3814,
			wantErr:      false,
		},
		{
			name: "Case error - wrong DistanceMeters",
			fields: fields{
				config: Config{
					fareRule: defaultFareRule,
				},
			},
			args: args{
				distanceMeters: []fareevaluator.DistanceMeter{
					{
						Mileage:     0,
						ElapsedTime: "00:00:00.000",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				config: tt.fields.config,
			}
			gotTaxiFare, err := s.CalculateFare(tt.args.distanceMeters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.CalculateFare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTaxiFare != tt.wantTaxiFare {
				t.Errorf("Service.CalculateFare() = %v, want %v", gotTaxiFare, tt.wantTaxiFare)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		options []Option
	}
	tests := []struct {
		name string
		args args
		want *Service
	}{
		{
			name: "Normal case",
			args: args{
				options: []Option{WithConfig(map[Scope]FareRule{})},
			},
			want: &Service{
				config: Config{
					fareRule: make(map[Scope]FareRule),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.options...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
