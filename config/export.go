package config

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"
)

// ENUM(yml,json,env)
type ExportType uint8

func Export(w io.Writer, cfg any, exportType ExportType, envPrefix string) error {
	switch exportType {
	case ExportTypeJson:
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(cfg)
	case ExportTypeYml:
		return yaml.NewEncoder(w).Encode(cfg)
	case ExportTypeEnv:
		return exportEnv(w, cfg, envPrefix, "")
	default:
		return fmt.Errorf("unsupported export type: %v", exportType)
	}
}

func exportEnv(w io.Writer, iface any, envPrefix string, prefix string) error {
	var sb strings.Builder
	if err := walkStruct(iface, prefix, func(fieldv reflect.Value, path, envKey string) error {

		key := envPrefix + envKey

		value := fmt.Sprintf("%v", fieldv.Interface())
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
		sb.WriteString("\n")

		return nil
	}); err != nil {
		return err
	}
	_, err := w.Write([]byte(sb.String()))
	return err
}
