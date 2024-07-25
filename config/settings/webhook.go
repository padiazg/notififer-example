package settings

import (
	"encoding/json"
	"regexp"
	"strings"
)

type WebhookSettings struct {
	Enabled      bool         `json:"enabled" yaml:"enabled"`
	Port         uint         `json:"port" yaml:"port"`
	UseTLS       bool         `json:"use-tls" yaml:"use-tls"`
	Certificates Certificates `json:"certificates" yaml:"certificates"`
}

type Certificate struct {
	Name     string //`json:"name" yaml:"name"`
	CertFile string //`json:"cert-file" yaml:"cert-file"`
	KeyFile  string //`json:"key-file" yaml:"key-file"`
}

type Certificates []Certificate

// UnmarshalText implements the encoding.TextUnmarshaler interface for the
// Certificates type.
// For it to work, the viper.Unmarshal() function must be called with the
// viper.DecodeHook() function to include mapstructure.TextUnmarshallerHookFunc() to
// the hooks list.
/*
viper.Unmarshal(s, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
	mapstructure.StringToTimeDurationHookFunc(),
	mapstructure.TextUnmarshallerHookFunc(),
)))
*/
func (c *Certificates) UnmarshalText(text []byte) error {
	// fmt.Printf("Certificates.UnmarshalText: %s\n", string(text))

	var (
		data         = strings.TrimSpace(string(text))
		re           = regexp.MustCompile(`(?m)(\[\s*\])`) // match empty array
		certificates = []Certificate{}
	)

	if data != "" && !re.MatchString(data) {
		if err := json.Unmarshal([]byte(text), &certificates); err != nil {
			return err
		}

		for _, certificate := range certificates {
			*c = append(*c, certificate)
		}
	}

	return nil
}

func (w *WebhookSettings) Valid() bool {
	return w.Port > 1024
}
