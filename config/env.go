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
func BindEnvs(v *viper.Viper, envPrefix string, iface any) { bindEnvs(v, envPrefix, "", iface) }

func bindEnvs(v *viper.Viper, envPrefix string, prefix string, iface any) {
	ifv := reflect.Indirect(reflect.ValueOf(iface))
	ift := ifv.Type()

	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := t.Name
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			name = tag
		}

		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		switch fieldv.Kind() {
		case reflect.Struct:
			bindEnvs(v, envPrefix, path, fieldv.Addr().Interface())
		default:
			envKey := strings.ToUpper(strings.ReplaceAll(path, ".", "_"))
			if err := v.BindEnv(path, envPrefix+envKey); err != nil {
				panic(err)
			}
		}
	}
}
