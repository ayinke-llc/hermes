package config

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type testConfig struct {
	Database struct {
		Host string `mapstructure:"host" json:"host"`
		Port int    `mapstructure:"port" json:"port"`
	} `mapstructure:"database" json:"database"`
	API struct {
		Key    string `mapstructure:"key" json:"key"`
		Secret string `mapstructure:"secret" json:"secret"`
	} `mapstructure:"api" json:"api"`
	Debug bool `mapstructure:"debug" json:"debug"`
}

func TestExport(t *testing.T) {
	cfg := testConfig{}
	cfg.Database.Host = "localhost"
	cfg.Database.Port = 5432
	cfg.API.Key = "test-key"
	cfg.API.Secret = "secret123"
	cfg.Debug = true

	envPrefix := "APP_"

	tests := []struct {
		name         string
		exportType   ExportType
		expectedFile string
	}{
		{
			name:         "json export",
			exportType:   ExportTypeJson,
			expectedFile: "expected.json",
		},
		{
			name:         "yml export",
			exportType:   ExportTypeYml,
			expectedFile: "expected.yml",
		},
		{
			name:         "env export",
			exportType:   ExportTypeEnv,
			expectedFile: "expected.env",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := Export(&buf, &cfg, tt.exportType, envPrefix)
			require.NoError(t, err)

			expectedPath := filepath.Join("testdata", tt.expectedFile)
			expected, err := os.ReadFile(expectedPath)
			require.NoError(t, err)

			require.Equal(t, string(expected), buf.String())
		})
	}
}
