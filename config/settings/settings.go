package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Settings struct {
	Webhook       WebhookSettings
	AMQP          AMQPSettings
	showValues    bool
	showKeyValues bool
}

func (s *Settings) Read(file string) error {
	if file == "" {
		viper.SetConfigFile(".notifier-example.yaml")
		viper.SetConfigType("yaml")
	} else {
		viper.SetConfigFile(file)
	}

	// fmt.Printf("ConfigFileUsed: %s\n", viper.ConfigFileUsed())

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("notifier")

	s.Defaults()
	viper.ReadInConfig()
	bindEnvs(s)

	if err := viper.Unmarshal(s, viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(
		mapstructure.StringToTimeDurationHookFunc(),
		// Needed to decode Webbook.Certificates. The Certificates type implements the UnmarshalText method.
		mapstructure.TextUnmarshallerHookFunc(),
		// Uncoment the next line to use the certificatesJsonMapHookFunc() hook for the Webhook.Certificates type. It's
		// here for example purposes. It must be placed before mapstructure.StringToSliceHookFunc(",")
		// certificatesJsonMapHookFunc(),
		mapstructure.StringToSliceHookFunc(","),
	))); err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	return nil
}

func (s *Settings) Defaults() {
	// viper.SetDefault("webserver.port", 7000)
	// viper.SetDefault("webserver.use-tls", false)
}

// Show shows the current config and the key/value pairs resulting from
// reading the config file and the env-vars
func (s *Settings) Show() {
	if s.showValues {
		b, err := json.MarshalIndent(s, "", "  ")
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Printf("Settings:\n%s\n", string(b))
	}

	if s.showKeyValues {
		fmt.Printf("Key/Value pairs:\n")
		// print all the keys
		keys := viper.AllKeys()
		slices.Sort(keys)
		for _, key := range keys {
			val := viper.Get(key)
			fmt.Printf("  %s: %v\n", key, val)
		}
	}
}

// Save saves the current config to a file
func (s *Settings) Save(name string) error {
	if err := viper.WriteConfigAs(name); err != nil {
		return fmt.Errorf("writing config file: %+v", err)
	}

	return nil
}

// SaveExample saves an example for the config file
func (s *Settings) SaveExample(name string) error {
	data, err := yaml.Marshal(&Settings{
		Webhook: WebhookSettings{
			Port:   7443,
			UseTLS: true,
			Certificates: Certificates{{
				Name:     "certificate-1",
				CertFile: "cert-file.crt",
				KeyFile:  "key-file.crt",
			}},
		},
		AMQP: AMQPSettings{
			Address: "amqp://localhost",
			Queue:   "notification",
		},
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(name, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ParseServeArgsAndFlags parses args and flags passed by cobra and stores it into
// the settings structure
func (s *Settings) ParseServeArgsAndFlags(cmd *cobra.Command, args []string) {
	s.showValues, _ = cmd.Flags().GetBool("show-config")
	s.showKeyValues, _ = cmd.Flags().GetBool("show-key-values")
}
