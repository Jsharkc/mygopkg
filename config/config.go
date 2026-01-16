package config

import (
	"errors"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

func Init(configName, fileType string, fpath []string, rawVal any) (*viper.Viper, error) {
	viperInst := viper.New()
	viperInst.SetConfigName(configName)
	viperInst.SetConfigType(fileType)
	for _, path := range fpath {
		viperInst.AddConfigPath(path)
	}
	err := viperInst.ReadInConfig()
	if err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	if err := viperInst.Unmarshal(rawVal); err != nil {
		return nil, err
	}

	log.Println("viper.AllSettings(): ", viperInst.AllSettings())
	return viperInst, nil
}

func IniWithEnv(configName, fileType string, fpath []string, rawVal any, parts ...string) (*viper.Viper, error) {
	viperInst := viper.New()

	// Set up environment variable handling first
	viperInst.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	prefix := strings.Join(parts, "_")
	if prefix != "" {
		viperInst.SetEnvPrefix(prefix)
	}

	// Set default values for all struct fields so AutomaticEnv binds them
	setDefaultsForEnv(viperInst, rawVal)

	viperInst.AutomaticEnv()

	viperInst.SetConfigName(configName) // name of config file (without extension)
	viperInst.SetConfigType(fileType)   // REQUIRED if the config file does not have the extension in the name
	for _, path := range fpath {
		viperInst.AddConfigPath(path)
	}

	// Attempt to read the config file
	if err := viperInst.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return nil, err
		}
	}

	if err := viperInst.Unmarshal(rawVal); err != nil {
		return nil, err
	}

	log.Println("viper.AllSettings(): ", viperInst.AllSettings())
	return viperInst, nil
}

// setDefaultsForEnv sets default values for all struct fields so AutomaticEnv will bind them
func setDefaultsForEnv(v *viper.Viper, iface interface{}) {
	ift := reflect.TypeOf(iface)

	if ift.Kind() == reflect.Ptr {
		ift = ift.Elem()
	}

	setDefaultsRecursive(v, ift, "")
}

func setDefaultsRecursive(v *viper.Viper, t reflect.Type, prefix string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagValue := field.Tag.Get("mapstructure")
		if tagValue == "" {
			tagValue = strings.ToLower(field.Name)
		}
		// Handle tag options like "name,omitempty"
		if idx := strings.Index(tagValue, ","); idx != -1 {
			tagValue = tagValue[:idx]
		}

		fullKey := tagValue
		if prefix != "" {
			fullKey = prefix + "." + tagValue
		}

		switch field.Type.Kind() {
		case reflect.Struct:
			setDefaultsRecursive(v, field.Type, fullKey)
		default:
			// Set default value so AutomaticEnv will bind this key to environment variable
			if !v.IsSet(fullKey) {
				v.SetDefault(fullKey, reflect.Zero(field.Type).Interface())
			}
		}
	}
}
