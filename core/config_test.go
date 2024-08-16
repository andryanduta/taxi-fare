package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"gopkg.in/gcfg.v1"
)

func TestInitMainConfig(t *testing.T) {
	originalReadFile := readFile
	defer func() { readFile = originalReadFile }() // Restore the original after the test

	tests := []struct {
		name    string
		setup   func()
		want    *Config
		wantErr bool
	}{
		{
			name: "Success case",
			setup: func() {
				// Replace readFile with a mock implementation
				readFile = func(cfg interface{}, filename string) error {
					tmpDir := t.TempDir()
					configFilePath := filepath.Join(tmpDir, "config-main.ini")

					configData := `[FaseRule "Over"]
Price = 40
Distance = 350
DistanceThreshold = 10000`
					if err := os.WriteFile(configFilePath, []byte(configData), 0644); err != nil {
						return err
					}

					return gcfg.ReadFileInto(cfg, configFilePath)
				}
			},
			want: &Config{
				FaseRule: map[string]*struct {
					Price             float64
					Distance          float64
					DistanceThreshold float64
				}{
					"Over": &struct {
						Price             float64
						Distance          float64
						DistanceThreshold float64
					}{
						Price:             40,
						Distance:          350,
						DistanceThreshold: 10000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "File read error",
			setup: func() {
				// Replace readFile with a mock implementation that returns an error
				readFile = func(cfg interface{}, filename string) error {
					return fmt.Errorf("file read error")
				}
			},
			want:    &Config{},
			wantErr: true,
		},
		{
			name: "Invalid file path",
			setup: func() {
				// Replace readFile with a mock implementation that simulates an invalid file path
				readFile = func(cfg interface{}, filename string) error {
					return gcfg.ReadFileInto(cfg, "/invalid/path/config-main.ini")
				}
			},
			want:    &Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got, err := InitMainConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitMainConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !compareConfig(got, tt.want) {
				t.Errorf("InitMainConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func compareConfig(got, want *Config) bool {
	if len(got.FaseRule) != len(want.FaseRule) {
		return false
	}

	for key, wantVal := range want.FaseRule {
		gotVal, ok := got.FaseRule[key]
		if !ok || !reflect.DeepEqual(gotVal, wantVal) {
			log.Println(gotVal, wantVal)
			return false
		}
	}

	return true
}
