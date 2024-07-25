package settings

type AMQPSettings struct {
	Enabled bool   `json:"enabled" yaml:"enabled"`
	Address string `json:"address" yaml:"address"`
	Queue   string `json:"queue" yaml:"queue"`
}
