package eviper

import (
	"github.com/spf13/viper"
	"os"
	"reflect"
	"strings"
)

type EViper struct {
	*viper.Viper
}

func New(v *viper.Viper) *EViper {
	return &EViper{v}
}

func (e *EViper) Unmarshal(rawVal interface{}) error {
	if err := e.Viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			// 	do nothing
		default:
			return err
		}
	}
	e.readEnvs(rawVal)
	if err := e.Viper.Unmarshal(rawVal); err != nil {
		return err
	}

	return nil
}

func (e *EViper) readEnvs(rawVal interface{}) {
	e.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	e.bindEnvs(rawVal, envs())
}

func (e *EViper) bindEnvs(in interface{}, envs map[string]string, prev ...string) {
	ifv := reflect.ValueOf(in)
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
	}

	for i := 0; i < ifv.NumField(); i++ {
		fv := ifv.Field(i)
		t := ifv.Type().Field(i)
		tv, hasEnvTag := t.Tag.Lookup("env")
		if hasEnvTag {
			if tv == ",squash" {
				e.bindEnvs(fv.Interface(), envs, prev...)
				continue
			}
		}
		name := t.Name
		mapTag, hasMapTag := t.Tag.Lookup("mapstructure")
		if hasMapTag {
			name = mapTag
		}
		env := strings.Join(append(prev, name), ".")
		switch fv.Kind() {
		case reflect.Struct:
			e.bindEnvs(fv.Interface(), envs, append(prev, t.Name)...)
		case reflect.Map:
			iter := fv.MapRange()
			for iter.Next() {
				if key, ok := iter.Key().Interface().(string); ok {
					e.bindEnvs(iter.Value().Interface(), envs, append(prev, t.Name, key)...)
				}
			}
		case reflect.Slice:
			e.Viper.SetTypeByDefaultValue(true)
			e.Viper.SetDefault(env, []string{})
			if hasEnvTag {
				if val, ok := envs[tv]; ok {
					splits := strings.Split(val, ",")
					e.Viper.Set(env, splits)
				}
			}
		default:
			if hasEnvTag {
				if val, ok := envs[tv]; ok {
					e.Viper.Set(env, val)
				}
			}
		}
	}
}

func envs() map[string]string {
	items := make(map[string]string)
	for _, item := range os.Environ() {
		splits := strings.Split(item, "=")
		items[splits[0]] = splits[1]
	}
	return items
}
