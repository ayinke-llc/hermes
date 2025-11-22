package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// BindEnvs takes the mapstructure tags of a struct and converts them to envs binding them to Viper.
//
//	type config struct {
//	  DSN string `mapstructure:"dns"`
//	}
//
// if your envPrefix is XXX, it can be set in env like XXX_DSN
func BindEnvs(v *viper.Viper, envPrefix string, iface any) {
	walkStruct(iface, "", func(fieldv reflect.Value, path, envKey string) error {
		if err := v.BindEnv(path, envPrefix+envKey); err != nil {
			panic(err)
		}
		return nil
	})
}

func walkStruct(iface any, prefix string, fn func(fieldv reflect.Value, path, envKey string) error) error {
	ifv := reflect.Indirect(reflect.ValueOf(iface))
	ift := ifv.Type()

	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)

		if !t.IsExported() {
			continue
		}

		if !fieldv.CanInterface() {
			continue
		}

		name := t.Name
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			if tag == "-" {
				continue
			}
			name = tag
		}

		jsonTag := t.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		envKey := strings.ToUpper(strings.ReplaceAll(path, ".", "_"))

		if fieldv.Kind() == reflect.Struct {
			if fieldv.CanAddr() && fieldv.Addr().CanInterface() {
				if err := walkStruct(fieldv.Addr().Interface(), path, fn); err != nil {
					return err
				}
			}
		} else {
			if err := fn(fieldv, path, envKey); err != nil {
				return err
			}
		}
	}
	return nil
}
