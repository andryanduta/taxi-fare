package main

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestSetupLogger(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		setup   func()
		want    bool
		wantErr bool
	}{
		{
			name: "Logger setup success",
			setup: func() {
				logFilePath := "logfile.log"
				if _, err := os.Stat(logFilePath); !os.IsNotExist(err) {
					os.Remove(logFilePath)
				}

				setupLogger()
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Log file exists",
			setup: func() {
				logFilePath := "logfile.log"
				if _, err := os.Stat(logFilePath); !os.IsNotExist(err) {
					os.Remove(logFilePath)
				}
				os.Create(logFilePath)

				setupLogger()
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			logFilePath := "logfile.log"

			var fileExists bool
			for i := 0; i < 5; i++ {
				if _, err := os.Stat(logFilePath); err == nil {
					fileExists = true
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			if fileExists != tt.want {
				if tt.want {
					t.Errorf("Log file was not created")
				} else {
					t.Errorf("Log file should not be created, but was found")
				}
			}

			// Clean up the log file after test
			if err := os.Remove(logFilePath); err != nil {
				t.Errorf("Failed to remove log file: %v", err)
			}

			if log.Logger.GetLevel() != zerolog.InfoLevel {
				t.Errorf("Expected log level to be InfoLevel, got %v", log.Logger.GetLevel())
			}
		})
	}
}
