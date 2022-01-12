package handler

import (
	"reflect"
	"testing"

	"github.com/andryanduta/taxi-fare/fareevaluator"
)

func Test_validateInputs(t *testing.T) {
	type args struct {
		inputs []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Case success",
			args: args{
				inputs: []string{
					"00:00:00.000 0", "00:03:00.123 700",
				},
			},
			wantErr: false,
		},
		{
			name: "Case error, invalid size",
			args: args{
				inputs: []string{
					"00:00:00.000 0",
				},
			},
			wantErr: true,
		},
		{
			name: "Case error, invalid format pattern",
			args: args{
				inputs: []string{
					"00:00:00.000 0", "00:03:00.123 700", "00:03:003 700",
				},
			},
			wantErr: true,
		},
		{
			name: "Case error, invalid mileage",
			args: args{
				inputs: []string{
					"00:00:00.000 0", "00:03:00.123 0",
				},
			},
			wantErr: true,
		},
		{
			name: "Case error, invalid elapsed time",
			args: args{
				inputs: []string{
					"00:00:00.000 0", "00:09:00.123 299",
				},
			},
			wantErr: true,
		},
		{
			name: "Case error, invalid elapsed time #2",
			args: args{
				inputs: []string{
					"00:00:00.000 0", "00:09:00.123 299", "00:09:00.123 299",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateInputs(tt.args.inputs); (err != nil) != tt.wantErr {
				t.Errorf("validateInputs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		svc fareevaluator.Service
	}
	tests := []struct {
		name string
		args args
		want *handler
	}{
		{
			name: "Success case",
			args: args{
				svc: nil,
			},
			want: &handler{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_constructInput(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name              string
		args              args
		wantDistanceMeter fareevaluator.DistanceMeter
		wantErr           bool
	}{
		{
			name: "Success case",
			args: args{
				line: "00:00:00.000 0",
			},
			wantDistanceMeter: fareevaluator.DistanceMeter{
				Mileage:     0,
				ElapsedTime: "00:00:00.000",
			},
		},
		{
			name: "Error case",
			args: args{
				line: "00:00:00.000 a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDistanceMeter, err := constructInput(tt.args.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("constructInput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDistanceMeter, tt.wantDistanceMeter) {
				t.Errorf("constructInput() = %v, want %v", gotDistanceMeter, tt.wantDistanceMeter)
			}
		})
	}
}
