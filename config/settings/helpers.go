package settings

import (
	"encoding/json"
	"log"
	"reflect"
	"regexp"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// bindEnvs creates the environment variable bindings for the given struct, also aliases for proper
// binding of environment variables and values from .env files and other structured config files.
func bindEnvs(i interface{}, parts ...string) {
	ifv := reflect.ValueOf(i)
	ift := reflect.TypeOf(i)

	// received a pointer, dereference it
	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for x := 0; x < ift.NumField(); x++ {
		// for _, f := range reflect.VisibleFields(ift) {
		t := ift.Field(x)
		v := ifv.Field(x)

		if !t.IsExported() {
			// fmt.Printf("  not exported: %s\n", t.Name)
			continue
		}

		// fmt.Printf("field: %s, value: %s\n", t.Name, v.String())
		switch v.Kind() {
		case reflect.Struct:
			// fmt.Printf("  struct: %s\n", t.Name)
			bindEnvs(v.Interface(), append(parts, t.Name)...)

		case reflect.Ptr:
			if v.IsNil() {
				// fmt.Printf("  nil pointer: %s\n", t.Name)
				continue
			}
			// fmt.Printf("  pointer: %s\n", t.Name)
			bindEnvs(v.Interface(), append(parts, t.Name)...)

		default:
			var (
				envKey   = strings.ToUpper(strings.Join(append(parts, t.Name), "_"))
				key      = strings.Join(append(parts, t.Name), ".")
				envAlias = strings.ToLower(envKey)
			)

			// set the env binding
			if err := viper.BindEnv(key, envKey); err != nil {
				log.Fatalf("config: unable to bind env: " + err.Error())
			}

			viper.RegisterAlias(envAlias, key)

			// fmt.Printf("  key: %s => %s => %s\n", key, envKey, envAlias)
		}
	}
}

// certificatesJsonMapHookFunc is a mapstructure.DecodeHookFunc that decodes a JSON string into a
// Certificates type.
// For it to work, the viper.Unmarshal() function must be called with the
// viper.DecodeHook() function to include it into the hooks list.
/*
viper.Unmarshal(s, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	mapstructure.StringToTimeDurationHookFunc(),
	certificatesJsonMapHookFunc(),
)))
*/
func certificatesJsonMapHookFunc() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
		// fmt.Printf("certificatesJsonMapHookFunc: \n\tf: %s, \n\tt: %s, \n\tdata: %s\n", f, t, data)
		if f.Kind() != reflect.String || t != reflect.TypeOf(Certificates{}) {
			return data, nil
		}

		var (
			text         = strings.TrimSpace(data.(string))
			re           = regexp.MustCompile(`(?m)(\[\s*\])`) // match empty array
			certificates = Certificates{}
			cert         = []Certificate{}
		)

		// decode only if the string is not an empty array
		if text != "" && !re.MatchString(text) {
			if err := json.Unmarshal([]byte(text), &cert); err != nil {
				return nil, err
			}

			for _, c := range cert {
				certificates = append(certificates, c)
			}
		}

		return certificates, nil
	}
}
